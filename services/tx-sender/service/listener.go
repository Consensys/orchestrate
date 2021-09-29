package service

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/pkg/utils/envelope"
	utils2 "github.com/consensys/orchestrate/services/tx-sender/tx-sender/utils"
	"google.golang.org/protobuf/proto"

	"github.com/Shopify/sarama"
	encoding "github.com/consensys/orchestrate/pkg/encoding/proto"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/tx"
	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"
)

const messageListenerComponent = "service.kafka-consumer"

type MessageListener struct {
	useCases     usecases.UseCases
	recoverTopic string
	crafterTopic string
	retryBackOff backoff.BackOff
	producer     sarama.SyncProducer
	jobClient    client.JobClient
	cancel       context.CancelFunc
	err          error
	logger       *log.Logger
}

func NewMessageListener(useCases usecases.UseCases,
	jobClient client.JobClient,
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
		jobClient:    jobClient,
		logger:       log.NewLogger().SetComponent(messageListenerComponent),
	}
}

func (listener *MessageListener) Setup(session sarama.ConsumerGroupSession) error {
	listener.logger.WithContext(session.Context()).
		WithField("kafka.generation_id", session.GenerationID()).
		WithField("kafka.member_id", session.MemberID()).
		WithField("claims", session.Claims()).
		Info("ready to consume messages")

	return nil
}

func (listener *MessageListener) Cleanup(session sarama.ConsumerGroupSession) error {
	logger := listener.logger.WithContext(session.Context())
	logger.Info("all claims consumed")
	if listener.cancel != nil {
		logger.Debug("canceling context")
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
	logger := listener.logger.WithContext(ctx)
	ctx = log.With(ctx, logger)
	logger.Info("started consuming claims loop")

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

			evlp, err := decodeMessage(logger, msg)
			if err != nil {
				logger.WithError(err).Error("error decoding message", msg)
				session.MarkMessage(msg, "")
				continue
			}

			logger.WithField("envelope_id", evlp.ID).
				WithField("timestamp", msg.Timestamp).
				Debug("message consumed")

			jlogger := logger.WithField("job", evlp.GetJobUUID()).WithField("schedule", evlp.GetScheduleUUID())
			err = listener.processEnvelope(log.With(ctx, jlogger), evlp)
			if err != nil {
				jlogger.WithError(err).Error("error processing message")
				return err
			}

			jlogger.Debug("job processed successfully")
			session.MarkMessage(msg, "")
		}
	}
}

func (listener *MessageListener) processEnvelope(ctx context.Context, evlp *tx.Envelope) error {
	logger := log.FromContext(ctx)
	tenantID := evlp.GetHeadersValue(utils.TenantIDMetadata)
	job := envelope.NewJobFromEnvelope(evlp, tenantID)
	return backoff.RetryNotify(
		func() error {
			err := listener.executeSendJob(ctx, job)
			switch {
			// Exits if not errors
			case err == nil:
				return nil
			case err == context.DeadlineExceeded || err == context.Canceled:
				return backoff.Permanent(err)
			case ctx.Err() != nil:
				return backoff.Permanent(ctx.Err())
			case errors.IsConnectionError(err):
				return err
			}

			var serr error
			switch {
			// Never retry on children jobs
			case job.InternalData.ParentJobUUID == job.UUID:
				serr = utils2.UpdateJobStatus(ctx, listener.jobClient, job,
					entities.StatusFailed, err.Error(), nil)
			// Retry over same message
			case errors.IsInvalidNonceWarning(err):
				resetEnvelopeTx(evlp)
				serr = utils2.UpdateJobStatus(ctx, listener.jobClient, job,
					entities.StatusRecovering, err.Error(), nil)
				if serr == nil {
					return err
				}
			// In case of other kind of errors...
			default:
				txResponse := evlp.AppendError(errors.FromError(err)).TxResponse()
				serr = listener.sendEnvelope(ctx, evlp.ID, txResponse, listener.recoverTopic, evlp.PartitionKey())
				if serr == nil {
					serr = utils2.UpdateJobStatus(ctx, listener.jobClient, job,
						entities.StatusFailed, err.Error(), nil)
				}
			}

			if serr != nil {
				// IMPORTANT: Jobs can be updated in parallel to NEVER_MINED or MINED, so that we should
				// ignore it in this case
				if strings.Contains(err.Error(), "42400@") {
					logger.WithError(err).Warn("ignored error")
					return nil
				}

				if ctx.Err() != nil {
					return backoff.Permanent(ctx.Err())
				}

				// Retry on IsConnectionError and not context canceled
				if errors.IsConnectionError(serr) {
					return serr
				}

				return backoff.Permanent(serr)
			}

			return nil
		},
		listener.retryBackOff,
		func(err error, duration time.Duration) {
			logger.WithError(err).Warnf("error processing job, retrying in %v...", duration)
		},
	)
}

func (listener *MessageListener) executeSendJob(ctx context.Context, job *entities.Job) error {
	switch string(job.Type) {
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

func decodeMessage(logger *log.Logger, msg *sarama.ConsumerMessage) (*tx.Envelope, error) {
	txEnvelope := &tx.TxEnvelope{}
	err := encoding.Unmarshal(msg.Value, txEnvelope)
	if err != nil {
		errMessage := "failed to decode request message"
		logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage).ExtendComponent(messageListenerComponent)
	}

	evlp, err := txEnvelope.Envelope()
	if err != nil {
		errMessage := "failed to extract envelope from request"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(messageListenerComponent)
	}

	return evlp, nil
}

func (listener *MessageListener) sendEnvelope(ctx context.Context, msgID string, protoMessage proto.Message, topic, partitionKey string) error {
	logger := listener.logger.WithContext(ctx).WithField("topic", topic).WithField("envelope_id", msgID)
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
		logger.WithError(err).Error(errMessage)
		return errors.EncodingError(errMessage).ExtendComponent(messageListenerComponent)
	}
	msg.Value = sarama.ByteEncoder(b)

	partition, offset, err := listener.producer.SendMessage(msg)
	if err != nil {
		errMessage := "failed to produce kafka message"
		logger.WithError(err).Error(errMessage)
		return errors.KafkaConnectionError(errMessage).ExtendComponent(messageListenerComponent)
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
