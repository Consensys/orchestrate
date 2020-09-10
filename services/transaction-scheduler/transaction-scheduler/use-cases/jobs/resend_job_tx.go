package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=resend_job_tx.go -destination=mocks/resend_job_tx.go -package=mocks

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
	if jobEntity.GetStatus() != utils.StatusPending {
		errMessage := "cannot resend job transaction at the current status"
		logger.WithField("status", jobEntity.GetStatus()).Error(errMessage)
		return errors.InvalidStateError(errMessage)
	}

	partition, offset, err := utils2.SendJobMessage(ctx, jobModel, uc.kafkaProducer, uc.topicsCfg.Sender)
	if err != nil {
		return errors.FromError(err).ExtendComponent(resendJobTxComponent)
	}

	logger.WithField("partition", partition).WithField("offset", offset).Info("job resend successfully")

	return nil
}
