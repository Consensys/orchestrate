package jobs

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils/envelope"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"

	pkgsarama "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/Shopify/sarama"
)

const resendJobTxComponent = "use-cases.resend-job-tx"

type resendJobTxUseCase struct {
	db            store.DB
	kafkaProducer sarama.SyncProducer
	topicsCfg     *pkgsarama.KafkaTopicConfig
	logger        *log.Logger
}

func NewResendJobTxUseCase(db store.DB, kafkaProducer sarama.SyncProducer, topicsCfg *pkgsarama.KafkaTopicConfig) usecases.ResendJobTxUseCase {
	return &resendJobTxUseCase{
		db:            db,
		kafkaProducer: kafkaProducer,
		topicsCfg:     topicsCfg,
		logger:        log.NewLogger().SetComponent(resendJobTxComponent),
	}
}

// Execute sends a job to the Kafka topic
func (uc *resendJobTxUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) error {
	ctx = log.WithFields(ctx, log.Field("job", jobUUID))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("resending job transaction")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants, false)
	if err != nil {
		return errors.FromError(err).ExtendComponent(resendJobTxComponent)
	}

	jobModel.InternalData.ParentJobUUID = jobUUID
	jobEntity := parsers.NewJobEntityFromModels(jobModel)
	if jobEntity.Status != entities.StatusPending {
		errMessage := "cannot resend job transaction at the current status"
		logger.WithField("status", jobEntity.Status).Error(errMessage)
		return errors.InvalidStateError(errMessage)
	}

	partition, offset, err := envelope.SendJobMessage(jobEntity, uc.kafkaProducer, uc.topicsCfg.Sender)
	if err != nil {
		logger.WithError(err).Error("failed to send job message")
		return errors.FromError(err).ExtendComponent(resendJobTxComponent)
	}

	logger.WithField("partition", partition).WithField("offset", offset).Info("job resend successfully")
	return nil
}
