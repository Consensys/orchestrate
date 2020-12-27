package service

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils/envelope"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/utils"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/use-cases"
)

const messageListenerComponent = "service.message-listener"

type MessageListener struct {
	useCases     usecases.UseCases
	recoverTopic string
	crafterTopic string
	retryBackOff backoff.BackOff
	producer     sarama.SyncProducer
	jobClient    client.JobClient
	cancel       context.CancelFunc
	err          error
}

func NewMessageListener(useCases usecases.UseCases, jobClient client.JobClient,
	producer sarama.SyncProducer, recoverTopic, crafterTopic string, bck backoff.BackOff) *MessageListener {
	return &MessageListener{
		useCases:     useCases,
		recoverTopic: recoverTopic,
		crafterTopic: crafterTopic,
		producer:     producer,
		retryBackOff: bck,
		jobClient:    jobClient,
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
	if listener.cancel != nil {
		log.WithContext(session.Context()).Debug("listener: canceling context")
		listener.cancel()
	}

	return listener.err
}

func (listener *MessageListener) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var ctx context.Context
	ctx, listener.cancel = context.WithCancel(session.Context())
	listener.err = listener.consumeClaimLoop(ctx, session, claim)
	return listener.err
}

func (listener *MessageListener) consumeClaimLoop(ctx context.Context, session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger := log.WithContext(ctx)
	logger.Info("tx-sender has started consuming claims loop")
	for {
		select {
		case <-ctx.Done():
			logger.WithField("reason", ctx.Err().Error()).Info("gracefully stopping message listener...")
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
					logger.WithField("job_uuid", job.UUID).Info("processing job")
					err = listener.processJob(ctx, job)
					switch {
					// Exits if not errors
					case err == nil:
						return nil
					case errors.IsConnectionError(err):
						return err
					case err == context.DeadlineExceeded || err == context.Canceled:
						return backoff.Permanent(err)
					case ctx.Err() != nil:
						return backoff.Permanent(ctx.Err())
					}

					var serr error
					switch {
					// Never retry on children jobs
					case job.InternalData.ParentJobUUID == job.UUID:
						serr = utils2.UpdateJobStatus(ctx, listener.jobClient, evlp.GetJobUUID(),
							utils.StatusFailed, err.Error(), nil)
					// Retry over same message
					case errors.IsInvalidNonceWarning(err):
						resetEnvelopeTx(evlp)
						serr = utils2.UpdateJobStatus(ctx, listener.jobClient, evlp.GetJobUUID(),
							utils.StatusRecovering, err.Error(), nil)
						if serr == nil {
							return err
						}
					// In case of other kind of errors...
					default:
						txResponse := evlp.AppendError(errors.FromError(err)).TxResponse()
						serr = listener.sendEnvelope(ctx, evlp.ID, txResponse, listener.recoverTopic, evlp.PartitionKey())
						if serr == nil {
							serr = utils2.UpdateJobStatus(ctx, listener.jobClient, evlp.GetJobUUID(),
								utils.StatusFailed, err.Error(), nil)
						}
					}

					if serr != nil {
						// IMPORTANT: Jobs can be updated in parallel to NEVER_MINED or MINED, so that we should
						// ignore it in this case
						if strings.Contains(err.Error(), "42400@") {
							logger.WithError(err).Warn("ignored error")
							return nil
						}

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

			logger.WithField("job_uuid", job.UUID).Info("job processed successfully")
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
