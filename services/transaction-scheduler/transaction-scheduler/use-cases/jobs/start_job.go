package jobs

import (
	"context"
	"fmt"

	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"

	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=start_job.go -destination=mocks/start_job.go -package=mocks

const startJobComponent = "use-cases.start-job"

type StartJobUseCase interface {
	Execute(ctx context.Context, jobUUID, tenantID string) error
}

// startJobUseCase is a use case to start a transaction job
type startJobUseCase struct {
	db             store.DB
	kafkaProducer  sarama.SyncProducer
	txCrafterTopic string
}

// NewStartJobUseCase creates a new StartJobUseCase
func NewStartJobUseCase(
	db store.DB,
	kafkaProducer sarama.SyncProducer,
	txCrafterTopic string,
) StartJobUseCase {
	return &startJobUseCase{
		db:             db,
		kafkaProducer:  kafkaProducer,
		txCrafterTopic: txCrafterTopic,
	}
}

// Execute validates and creates a new transaction job
func (uc *startJobUseCase) Execute(ctx context.Context, jobUUID, tenantID string) error {
	logger := log.WithContext(ctx)

	logger.
		WithField("job_uuid", jobUUID).
		Debugf("starting job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenantID)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	txEnvelope := parsers.NewEnvelopeFromJobModel(jobModel)
	partition, offset, err := uc.sendMessage(ctx, txEnvelope)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	logger.
		WithField("envelope_id", txEnvelope.GetID()).
		Debugf("envelope sent to kafka")

	jobLog := &models.Log{
		JobID:  &jobModel.ID,
		Status: entities.JobStatusStarted,
		Message: fmt.Sprintf("message sent to partition %v, offset %v and topic %v", partition, offset,
			uc.txCrafterTopic),
	}

	if err = uc.db.Log().Insert(ctx, jobLog); err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	logger.
		WithField("job_uuid", jobUUID).
		Info("job started successfully")

	return nil
}

func (uc *startJobUseCase) sendMessage(ctx context.Context, txEnvelope proto.Message) (partition int32, offset int64, err error) {
	log.WithContext(ctx).Debug("sending kafka message")

	envelopeBytes, err := encoding.Marshal(txEnvelope)
	if err != nil {
		errMessage := "failed to encode envelope"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.InvalidParameterError(errMessage)
	}

	msg := &sarama.ProducerMessage{
		Topic:   uc.txCrafterTopic,
		Value:   sarama.ByteEncoder(envelopeBytes),
		Headers: nil, // TODO: Add the JWT token here? https://pegasys1.atlassian.net/browse/PO-544
	}

	// Send message
	partition, offset, err = uc.kafkaProducer.SendMessage(msg)
	if err != nil {
		errMessage := "could not produce kafka message"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return 0, 0, errors.KafkaConnectionError(errMessage).ExtendComponent(startJobComponent)
	}

	return partition, offset, err
}
