package service

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer/tx-signer/use-cases"

	"github.com/Shopify/sarama"
)

const messageListenerComponent = "service.message-listener"
const recoverableErrorMessage = "retrying message on recoverable error"

type MessageListener struct {
	useCases          usecases.UseCases
	senderTopic       string
	recoverTopic      string
	txSchedulerClient client.TransactionSchedulerClient
}

func NewMessageListener(ucs usecases.UseCases, senderTopic, recoverTopic string, txSchedulerClient client.TransactionSchedulerClient) *MessageListener {
	return &MessageListener{
		useCases:          ucs,
		senderTopic:       senderTopic,
		recoverTopic:      recoverTopic,
		txSchedulerClient: txSchedulerClient,
	}
}

func (MessageListener) Setup(session sarama.ConsumerGroupSession) error {
	log.WithContext(session.Context()).
		WithField("kafka.generation_id", session.GenerationID()).
		WithField("kafka.member_id", session.MemberID()).
		WithField("claims", session.Claims()).
		Info("listener ready to consume messages")

	return nil
}

func (MessageListener) Cleanup(session sarama.ConsumerGroupSession) error {
	log.WithContext(session.Context()).Info("listener: all claims consumed")
	return nil
}

func (listener *MessageListener) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	logger := log.WithContext(ctx)
	logger.Info("tx-signer has started consuming claims")

	for {
		select {
		case msg := <-claim.Messages():
			envelope, err := decodeMessage(msg)
			if err != nil {
				session.MarkMessage(msg, "")
				continue
			}

			// Skip processing on raw transactions
			raw, txHash, err := listener.processEnvelope(ctx, envelope)

			var der error
			switch {
			case err != nil && errors.IsConnectionError(err):
				logger.Error(recoverableErrorMessage)
				continue
			case err != nil:
				txResponse := envelope.AppendError(errors.FromError(err)).TxResponse()
				der = listener.useCases.SendEnvelope().Execute(ctx, txResponse, listener.recoverTopic, envelope.PartitionKey())
			default:
				_ = envelope.SetRawString(raw)
				_ = envelope.SetTxHashString(txHash)
				der = listener.useCases.SendEnvelope().Execute(ctx, envelope.TxEnvelopeAsRequest(), listener.senderTopic, envelope.PartitionKey())
			}
			if der != nil && errors.IsConnectionError(der) {
				logger.Error(recoverableErrorMessage)
				continue
			}

			if der != nil {
				err = listener.updateTransactionStatus(ctx, envelope.GetJobUUID(), err.Error())
				if err != nil && errors.IsConnectionError(err) {
					logger.Error(recoverableErrorMessage)
					continue
				}
			}

			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			logger.WithField("reason", ctx.Err().Error()).Info("gracefully stopping message listener...")
			return nil
		}
	}
}

func decodeMessage(msg *sarama.ConsumerMessage) (*tx.Envelope, error) {
	txEnvelope := &tx.TxEnvelope{}
	err := proto.Unmarshal(msg.Value, txEnvelope)
	if err != nil {
		errMessage := "failed to decode request message"
		log.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage).ExtendComponent(messageListenerComponent)
	}

	envelope, err := txEnvelope.Envelope()
	if err != nil {
		errMessage := "failed to extract envelope from request"
		log.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(messageListenerComponent)
	}

	return envelope, nil
}

func (listener *MessageListener) processEnvelope(ctx context.Context, envelope *tx.Envelope) (raw, txHash string, err error) {
	tenantID := envelope.GetHeadersValue(multitenancy.TenantIDMetadata)
	job := EnvelopeToJob(envelope, tenantID)

	switch {
	case envelope.IsEthSendRawTransaction():
		log.WithContext(ctx).WithField("job_uuid", envelope.GetJobUUID()).Info("raw transaction processed successfully")
		return job.Transaction.Raw, "", nil
	case envelope.IsEthSendTesseraPrivateTransaction():
		log.WithContext(ctx).WithField("job_uuid", envelope.GetJobUUID()).Info("tessera transaction processed successfully")
		// Do nothing as we do not sign storeRaw payload
		return "", "", nil
	case envelope.IsEthSendTesseraMarkingTransaction():
		return listener.useCases.SignQuorumPrivateTransaction().Execute(ctx, job)
	case envelope.IsEeaSendPrivateTransaction():
		return listener.useCases.SignEEATransaction().Execute(ctx, job)
	default:
		return listener.useCases.SignTransaction().Execute(ctx, job)
	}
}

func (listener *MessageListener) updateTransactionStatus(ctx context.Context, jobUUID, errMessage string) error {
	_, err := listener.txSchedulerClient.UpdateJob(ctx, jobUUID, &txschedulertypes.UpdateJobRequest{
		Status:  utils.StatusFailed,
		Message: errMessage,
	})
	if err != nil {
		errMessage := "failed to update transaction status"
		log.WithError(err).WithField("status", utils.StatusFailed).Error(errMessage)
		return errors.FromError(err).ExtendComponent(messageListenerComponent)
	}

	return nil
}
