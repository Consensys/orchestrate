package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils/envelope"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const resendJobTxComponent = "use-cases.resend-job-tx"

// resendJobTxUseCase is a use case to start a transaction job
type resendJobTxUseCase struct {
	db            store.DB
	kafkaProducer sarama.SyncProducer
	topicsCfg     *pkgsarama.KafkaTopicConfig
}

// NewResendJobTxUseCase creates a new ResendJobTxUseCase
func NewResendJobTxUseCase(db store.DB, kafkaProducer sarama.SyncProducer, topicsCfg *pkgsarama.KafkaTopicConfig) usecases.ResendJobTxUseCase {
	return &resendJobTxUseCase{
		db:            db,
		kafkaProducer: kafkaProducer,
		topicsCfg:     topicsCfg,
	}
}

// Execute sends a job to the Kafka topic
func (uc *resendJobTxUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) error {
	logger := log.WithContext(ctx).WithField("job_uuid", jobUUID)
	logger.Debug("resending job transaction")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(resendJobTxComponent)
	}

	jobModel.InternalData.ParentJobUUID = jobUUID
	jobEntity := parsers.NewJobEntityFromModels(jobModel)
	if jobEntity.Status != utils.StatusPending {
		errMessage := "cannot resend job transaction at the current status"
		logger.WithField("status", jobEntity.Status).Error(errMessage)
		return errors.InvalidStateError(errMessage)
	}

	partition, offset, err := envelope.SendJobMessage(ctx, jobEntity, uc.kafkaProducer, uc.topicsCfg.Signer)
	if err != nil {
		return errors.FromError(err).ExtendComponent(resendJobTxComponent)
	}

	logger.WithField("partition", partition).WithField("offset", offset).Info("job resend successfully")

	return nil
}
