package service

import (
	"context"
	"time"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang/protobuf/proto"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils/envelope"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases"
)

const messageListenerComponent = "service.message-listener"

type MessageListener struct {
	useCases     usecases.UseCases
	recoverTopic string
	crafterTopic string
	retryBackOff backoff.BackOff
	producer     sarama.SyncProducer
	client       orchestrateclient.OrchestrateClient
	breakLoop    chan bool
}

func NewMessageListener(
	useCases usecases.UseCases,
	client orchestrateclient.OrchestrateClient,
	producer sarama.SyncProducer,
	recoverTopic, crafterTopic string,
	bck backoff.BackOff,
) *MessageListener {
	return &MessageListener{
		useCases:     useCases,
		recoverTopic: recoverTopic,
		crafterTopic: crafterTopic,
		producer:     producer,
		retryBackOff: bck,
		client:       client,
		breakLoop:    make(chan bool, 1),
	}
}

func (listener *MessageListener) Setup(session sarama.ConsumerGroupSession) error {
	log.WithContext(session.Context()).
		WithField("kafka.generation_id", session.GenerationID()).
		WithField("kafka.member_id", session.MemberID()).
		WithField("claims", session.Claims()).
		Info("listener ready to consume messages")

	return nil
}

func (listener *MessageListener) Cleanup(session sarama.ConsumerGroupSession) error {
	log.WithContext(session.Context()).Info("listener: all claims consumed")
	return nil
}

func (listener *MessageListener) Break(session sarama.ConsumerGroupSession) error {
	log.WithContext(session.Context()).Info("listener: has been stopped")
	listener.breakLoop <- true
	if session != nil {
		return listener.Cleanup(session)
	}
	return nil
}

func (listener *MessageListener) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	logger := log.WithContext(ctx)
	logger.Info("tx-signer has started consuming claims")

	for {
		select {
		case <-session.Context().Done():
			logger.WithField("reason", ctx.Err().Error()).Info("gracefully stopping message listener...")
			return nil
		case <-listener.breakLoop:
			logger.Info("gracefully stopping message listener...")
			return nil
		case msg, ok := <-claim.Messages():
			// Input channel has been close so we leave the loop
			if !ok {
				return nil
			}

			evlp, err := decodeMessage(msg)
			if err != nil {
				logger.WithError(err).Error("error decoding message", msg)
				session.MarkMessage(msg, "")
				continue
			}

			logger.WithField("envelope_id", evlp.ID).WithField("timestamp", msg.Timestamp).
				Info("message consumed")

			tenantID := evlp.GetHeadersValue(utils.TenantIDMetadata)
			job := envelope.NewJobFromEnvelope(evlp, tenantID)

			err = backoff.RetryNotify(
				func() error {
					err = listener.processJob(ctx, job)
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

					var serr error
					if errors.IsInvalidNonceWarning(err) {
						resetEnvelopeTx(evlp)
						_, _, serr = envelope.SendJobMessage(ctx, job, listener.producer, listener.crafterTopic)
						if serr == nil {
							_, err = listener.client.UpdateJob(ctx, evlp.GetJobUUID(), &types.UpdateJobRequest{
								Status:  utils.StatusRecovering,
								Message: err.Error(),
							})
						}
					} else {
						// In case of other kind of errors...
						txResponse := evlp.AppendError(errors.FromError(err)).TxResponse()
						serr = listener.sendEnvelope(ctx, evlp.ID, txResponse, listener.recoverTopic, evlp.PartitionKey())
						if serr == nil {
							_, err = listener.client.UpdateJob(ctx, evlp.GetJobUUID(), &types.UpdateJobRequest{
								Status:  utils.StatusFailed,
								Message: err.Error(),
							})
						}
					}

					if serr != nil {
						// Retry on IsConnectionError
						if errors.IsConnectionError(serr) {
							return serr
						}
						return backoff.Permanent(serr)
					}

					return nil
				},
				listener.retryBackOff,
				func(err error, duration time.Duration) {
					logger.WithError(err).WithField("job", job.UUID).Warnf("error processing job, retrying in %v...", duration)
				},
			)

			if err != nil {
				logger.WithError(err).Errorf("error processing message")
				return err
			}

			session.MarkMessage(msg, "")
		}
	}
}

func (listener *MessageListener) processJob(ctx context.Context, job *entities.Job) error {
	switch job.Type {
	case tx.JobType_ETH_TESSERA_PRIVATE_TX.String():
		return listener.useCases.SendTesseraPrivateTx().Execute(ctx, job)
	case tx.JobType_ETH_TESSERA_MARKING_TX.String():
		return listener.useCases.SendTesseraMarkingTx().Execute(ctx, job)
	case tx.JobType_ETH_ORION_EEA_TX.String():
		return listener.useCases.SendEEAPrivateTx().Execute(ctx, job)
	case tx.JobType_ETH_RAW_TX.String():
		return listener.useCases.SendETHRawTx().Execute(ctx, job)
	case tx.JobType_ETH_ORION_MARKING_TX.String(), tx.JobType_ETH_TX.String():
		return listener.useCases.SendETHTx().Execute(ctx, job)
	default:
		return errors.InvalidParameterError("job type %s is not supported", job.Type)
	}
}

func decodeMessage(msg *sarama.ConsumerMessage) (*tx.Envelope, error) {
	txEnvelope := &tx.TxEnvelope{}
	err := encoding.Unmarshal(msg.Value, txEnvelope)
	if err != nil {
		errMessage := "failed to decode request message"
		log.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage).ExtendComponent(messageListenerComponent)
	}

	evlp, err := txEnvelope.Envelope()
	if err != nil {
		errMessage := "failed to extract envelope from request"
		log.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(messageListenerComponent)
	}

	return evlp, nil
}

func (listener *MessageListener) sendEnvelope(ctx context.Context, msgID string, protoMessage proto.Message, topic, partitionKey string) error {
	logger := log.WithContext(ctx).WithField("topic", topic).WithField("envelope_id", msgID)
	logger.Debug("sending envelope")

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	// Set key for Kafka partitions
	if partitionKey != "" {
		msg.Key = sarama.StringEncoder(partitionKey)
	}

	b, err := encoding.Marshal(protoMessage)
	if err != nil {
		errMessage := "failed to marshal envelope as request"
		log.WithError(err).Error(errMessage)
		return errors.EncodingError(errMessage)
	}
	msg.Value = sarama.ByteEncoder(b)

	partition, offset, err := listener.producer.SendMessage(msg)
	if err != nil {
		errMessage := "failed to produce kafka message"
		log.WithError(err).Error(errMessage)
		return errors.KafkaConnectionError(errMessage)
	}

	logger.WithField("partition", partition).
		WithField("offset", offset).
		Info("envelope successfully sent")

	return nil
}

func resetEnvelopeTx(req *tx.Envelope) {
	req.Nonce = nil
	req.TxHash = nil
	req.Raw = ""
}
