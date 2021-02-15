package jobs

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils/envelope"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"

	"github.com/Shopify/sarama"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const startJobComponent = "use-cases.start-job"

// startJobUseCase is a use case to start a transaction job
type startJobUseCase struct {
	db            store.DB
	kafkaProducer sarama.SyncProducer
	topicsCfg     *pkgsarama.KafkaTopicConfig
	metrics       metrics.TransactionSchedulerMetrics
	logger        *log.Logger
}

// NewStartJobUseCase creates a new StartJobUseCase
func NewStartJobUseCase(
	db store.DB,
	kafkaProducer sarama.SyncProducer,
	topicsCfg *pkgsarama.KafkaTopicConfig,
	m metrics.TransactionSchedulerMetrics,
) usecases.StartJobUseCase {
	return &startJobUseCase{
		db:            db,
		kafkaProducer: kafkaProducer,
		topicsCfg:     topicsCfg,
		metrics:       m,
		logger:        log.NewLogger().SetComponent(startJobComponent),
	}
}

// Execute sends a job to the Kafka topic
func (uc *startJobUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) error {
	logger := uc.logger.WithContext(ctx).WithField("job", jobUUID)
	logger.Debug("starting job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	jobEntity := parsers.NewJobEntityFromModels(jobModel)
	if !canUpdateStatus(entities.StatusStarted, jobEntity.Status) {
		errMessage := "cannot start job at the current status"
		logger.WithField("status", jobEntity.Status).WithField("next_status", entities.StatusStarted).Error(errMessage)
		return errors.InvalidStateError(errMessage)
	}

	jobModel.Status = entities.StatusStarted
	jobLog := &models.Log{
		JobID:  &jobModel.ID,
		Status: entities.StatusStarted,
	}

	dbtx, err := uc.db.Begin()
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	if err = dbtx.(store.Tx).Job().Update(ctx, jobModel); err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	if err = dbtx.(store.Tx).Log().Insert(ctx, jobLog); err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	partition, offset, err := envelope.SendJobMessage(jobEntity, uc.kafkaProducer, uc.topicsCfg.Sender)
	if err != nil {
		logger.WithError(err).Error("failed to send job message")
		_ = dbtx.Rollback()
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	if err := dbtx.Commit(); err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	uc.addMetrics(jobLog, jobModel.Logs[len(jobModel.Logs)-1], jobModel.ChainUUID)
	logger.WithField("partition", partition).WithField("offset", offset).Info("job started successfully")
	return nil
}

func (uc *startJobUseCase) addMetrics(current, previous *models.Log, chainUUID string) {
	baseLabels := []string{
		"chain_uuid", chainUUID,
	}

	d := float64(current.CreatedAt.Sub(previous.CreatedAt).Nanoseconds()) / float64(time.Second)
	uc.metrics.JobsLatencyHistogram().With(append(baseLabels,
		"prev_status", string(previous.Status),
		"status", string(current.Status),
	)...).Observe(d)
}
