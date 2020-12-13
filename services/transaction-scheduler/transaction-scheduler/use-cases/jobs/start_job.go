package jobs

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/metrics"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/use-cases"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/parsers"
)

const startJobComponent = "use-cases.start-job"

// startJobUseCase is a use case to start a transaction job
type startJobUseCase struct {
	db            store.DB
	kafkaProducer sarama.SyncProducer
	topicsCfg     *pkgsarama.KafkaTopicConfig
	metrics       metrics.TransactionSchedulerMetrics
}

// NewStartJobUseCase creates a new StartJobUseCase
func NewStartJobUseCase(db store.DB, kafkaProducer sarama.SyncProducer, topicsCfg *pkgsarama.KafkaTopicConfig, m metrics.TransactionSchedulerMetrics) usecases.StartJobUseCase {
	return &startJobUseCase{
		db:            db,
		kafkaProducer: kafkaProducer,
		topicsCfg:     topicsCfg,
		metrics:       m,
	}
}

// Execute sends a job to the Kafka topic
func (uc *startJobUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) error {
	logger := log.WithContext(ctx).WithField("job_uuid", jobUUID)
	logger.Debug("starting job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	jobEntity := parsers.NewJobEntityFromModels(jobModel)
	if !canUpdateStatus(utils.StatusStarted, jobEntity.Status) {
		errMessage := "cannot start job at the current status"
		logger.WithField("status", jobEntity.Status).WithField("next_status", utils.StatusStarted).Error(errMessage)
		return errors.InvalidStateError(errMessage)
	}

	var msgTopic string
	switch {
	case jobModel.Type == utils.EthereumRawTransaction:
		msgTopic = uc.topicsCfg.Signer
	default:
		msgTopic = uc.topicsCfg.Crafter
	}

	jobLog := &models.Log{
		JobID:  &jobModel.ID,
		Status: utils.StatusStarted,
	}

	dbtx, err := uc.db.Begin()
	if err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	if err = dbtx.(store.Tx).Log().Insert(ctx, jobLog); err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	partition, offset, err := envelope.SendJobMessage(ctx, jobEntity, uc.kafkaProducer, msgTopic)
	if err != nil {
		_ = dbtx.Rollback()
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	if err := dbtx.Commit(); err != nil {
		return errors.FromError(err).ExtendComponent(startJobComponent)
	}

	logger.WithField("partition", partition).WithField("offset", offset).Info("job started successfully")
	uc.addMetrics(jobLog, jobModel.Logs[len(jobModel.Logs)-1], jobModel.ChainUUID)

	return nil
}

func (uc *startJobUseCase) addMetrics(current, previous *models.Log, chainUUID string) {
	baseLabels := []string{
		"chain_uuid", chainUUID,
	}

	d := float64(current.CreatedAt.Sub(previous.CreatedAt).Nanoseconds()) / float64(time.Second)
	uc.metrics.JobsLatencyHistogram().With(append(baseLabels,
		"prev_status", previous.Status,
		"status", current.Status,
	)...).Observe(d)
}
