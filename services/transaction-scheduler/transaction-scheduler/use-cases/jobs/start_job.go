package jobs

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=start_job.go -destination=mocks/start_job.go -package=mocks

const startJobComponent = "use-cases.start-job"

type StartJobUseCase interface {
	Execute(ctx context.Context, jobUUID, tenantID string) error
}

// startJobUseCase is a use case to start a transaction job
type startJobUseCase struct {
	db            store.DB
	kafkaProducer sarama.SyncProducer
	topicsCfg     *pkgsarama.KafkaTopicConfig
}

// NewStartJobUseCase creates a new StartJobUseCase
func NewStartJobUseCase(db store.DB, kafkaProducer sarama.SyncProducer, topicsCfg *pkgsarama.KafkaTopicConfig) StartJobUseCase {
	return &startJobUseCase{
		db:            db,
		kafkaProducer: kafkaProducer,
		topicsCfg:     topicsCfg,
	}
}

// Execute sends a job to the Kafka topic
func (uc *startJobUseCase) Execute(ctx context.Context, jobUUID, tenantID string) error {
	logger := log.WithContext(ctx)

	logger.
		WithField("job_uuid", jobUUID).
		Debugf("starting job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenantID)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	var msgTopic string
	switch {
	case jobModel.Type == types.EthereumRawTransaction:
		msgTopic = uc.topicsCfg.Sender
	default:
		msgTopic = uc.topicsCfg.Crafter
	}

	partition, offset, err := uc.sendMessage(ctx, jobModel, msgTopic)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	jobLog := &models.Log{
		JobID:   &jobModel.ID,
		Status:  types.StatusStarted,
		Message: fmt.Sprintf("message sent to partition %v, offset %v and topic %v", partition, offset, msgTopic),
	}

	if err = uc.db.Log().Insert(ctx, jobLog); err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	logger.
		WithField("job_uuid", jobUUID).
		Info("job started successfully")

	return nil
}

func (uc *startJobUseCase) sendMessage(ctx context.Context, jobModel *models.Job, topic string) (partition int32, offset int64, err error) {
	log.WithContext(ctx).Debug("sending kafka message")

	txEnvelope := parsers.NewEnvelopeFromJobModel(jobModel, map[string]string{
		multitenancy.AuthorizationMetadata: authutils.AuthorizationFromContext(ctx),
	})

	evlp, err := txEnvelope.Envelope()
	if err != nil {
		errMessage := "failed to craft envelope"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.InvalidParameterError(errMessage)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(evlp.KafkaPartitionKey()),
	}

	err = encoding.Marshal(txEnvelope, msg)
	if err != nil {
		errMessage := "failed to encode envelope"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.InvalidParameterError(errMessage)
	}

	// Send message
	partition, offset, err = uc.kafkaProducer.SendMessage(msg)
	if err != nil {
		errMessage := "could not produce kafka message"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.KafkaConnectionError(errMessage).ExtendComponent(startJobComponent)
	}

	log.WithField("envelope_id", txEnvelope.GetID()).
		WithField("job_type", evlp.GetJobTypeString()).
		Debug("envelope sent to kafka")
	return partition, offset, err
}
