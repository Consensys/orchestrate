package service

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
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

type MessageListener struct {
	useCases          usecases.UseCases
	senderTopic       string
	recoverTopic      string
	txSchedulerClient client.TransactionSchedulerClient
	retryBackOff      backoff.BackOff
}

func NewMessageListener(ucs usecases.UseCases, senderTopic, recoverTopic string,
	txSchedulerClient client.TransactionSchedulerClient, bck backoff.BackOff) *MessageListener {
	return &MessageListener{
		useCases:          ucs,
		senderTopic:       senderTopic,
		recoverTopic:      recoverTopic,
		txSchedulerClient: txSchedulerClient,
		retryBackOff:      bck,
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
			if msg == nil {
				continue
			}

			envelope, err := decodeMessage(msg)
			if err != nil {
				logger.WithError(err).Error("error decoding message", msg)
				session.MarkMessage(msg, "")
				continue
			}

			err = backoff.RetryNotify(
				func() error {
					err = listener.processEnvelope(ctx, envelope)
					switch {
					// Exits if not errors
					case err == nil:
						return nil
					// Retry on IsConnectionError
					case errors.IsConnectionError(err):
						return err
					case err == context.DeadlineExceeded || err == context.Canceled:
						return backoff.Permanent(err)
					case ctx.Err() != nil:
						return backoff.Permanent(ctx.Err())
					}

					// In case of other kind of errors...
					txResponse := envelope.AppendError(errors.FromError(err)).TxResponse()
					err2 := listener.useCases.SendEnvelope().Execute(ctx, txResponse, listener.recoverTopic, envelope.PartitionKey())
					if err2 != nil {
						if errors.IsConnectionError(err2) {
							return err2
						}
						return backoff.Permanent(err2)
					}

					err = listener.updateTransactionStatus(ctx, envelope.GetJobUUID(), err.Error())
					if err != nil {
						if errors.IsConnectionError(err) {
							return err
						}
						return backoff.Permanent(err)
					}

					return nil
				},
				listener.retryBackOff,
				func(err error, duration time.Duration) {
					logger.WithError(err).Warnf("error processing envelope %q, retrying in %v...", envelope.ID, duration)
				},
			)

			if err != nil {
				logger.WithError(err).Error("error processing message", msg)
				return err
			}

			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			logger.WithField("reason", ctx.Err().Error()).Info("gracefully stopping message listener...")
			return nil
		}
	}
}

func (listener *MessageListener) processEnvelope(ctx context.Context, envelope *tx.Envelope) error {
	raw, txHash, err := listener.signEnvelopeTransaction(ctx, envelope)
	if err != nil {
		return err
	}

	_ = envelope.SetRawString(raw)
	_ = envelope.SetTxHashString(txHash)
	err = listener.useCases.SendEnvelope().Execute(ctx, envelope.TxEnvelopeAsRequest(), listener.senderTopic, envelope.PartitionKey())
	if err != nil {
		return err
	}

	return nil
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

func (listener *MessageListener) signEnvelopeTransaction(ctx context.Context, envelope *tx.Envelope) (raw, txHash string, err error) {
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
