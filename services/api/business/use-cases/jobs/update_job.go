package jobs

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const updateJobComponent = "use-cases.update-job"

// updateJobUseCase is a use case to create a new transaction job
type updateJobUseCase struct {
	db                    store.DB
	updateChildrenUseCase usecases.UpdateChildrenUseCase
	startNextJobUseCase   usecases.StartNextJobUseCase
	metrics               metrics.TransactionSchedulerMetrics
	logger                *log.Logger
}

// NewUpdateJobUseCase creates a new UpdateJobUseCase
func NewUpdateJobUseCase(db store.DB, updateChildrenUseCase usecases.UpdateChildrenUseCase,
	startJobUC usecases.StartNextJobUseCase, m metrics.TransactionSchedulerMetrics) usecases.UpdateJobUseCase {
	return &updateJobUseCase{
		db:                    db,
		updateChildrenUseCase: updateChildrenUseCase,
		startNextJobUseCase:   startJobUC,
		metrics:               m,
		logger:                log.NewLogger().SetComponent(updateJobComponent),
	}
}

// Execute validates and creates a new transaction job
func (uc *updateJobUseCase) Execute(ctx context.Context, job *entities.Job, nextStatus entities.JobStatus,
	logMessage string, tenants []string) (*entities.Job, error) {
	ctx = log.WithFields(ctx, log.Field("job", job.UUID), log.Field("next_status", nextStatus))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("updating job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants)
	if err != nil {
		return nil, err
	}

	// Does current job belong to a parent/children chains?
	var parentJobUUID string
	if jobModel.InternalData.ParentJobUUID != "" {
		parentJobUUID = jobModel.InternalData.ParentJobUUID
	} else if jobModel.InternalData.RetryInterval != 0 {
		parentJobUUID = job.UUID
	}

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		// We should lock ONLY when there is children jobs
		if parentJobUUID != "" {
			logger.WithField("parent_job", parentJobUUID).Debug("lock parent job row for update")
			der := tx.(store.Tx).Job().LockOneByUUID(ctx, parentJobUUID)
			if der != nil {
				return der
			}

			// Refresh jobModel after lock to ensure nothing was updated
			jobModel, err = uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants)
			if err != nil {
				return err
			}
		}

		if entities.IsFinalJobStatus(jobModel.Status) {
			errMessage := "job status is final, cannot be updated"
			logger.WithField("status", jobModel.Status).Error(errMessage)
			return errors.InvalidParameterError(errMessage).ExtendComponent(updateJobComponent)
		}

		// We are not forced to update the status
		if nextStatus != "" && !canUpdateStatus(nextStatus, jobModel.Status) {
			errMessage := "invalid status update for the current job state"
			logger.WithField("status", jobModel.Status).WithField("next_status", nextStatus).Error(errMessage)
			return errors.InvalidStateError(errMessage).ExtendComponent(updateJobComponent)
		}

		// We are not forced to update the transaction
		if job.Transaction != nil {
			parsers.UpdateTransactionModelFromEntities(jobModel.Transaction, job.Transaction)
			if der := tx.(store.Tx).Transaction().Update(ctx, jobModel.Transaction); der != nil {
				return der
			}
		}

		if nextStatus != "" {
			jobLogModel := &models.Log{
				JobID:   &jobModel.ID,
				Status:  nextStatus,
				Message: logMessage,
			}

			if der := tx.(store.Tx).Log().Insert(ctx, jobLogModel); der != nil {
				return der
			}

			jobModel.Logs = append(jobModel.Logs, jobLogModel)
		}

		updateJobModel(jobModel, job)
		if der := tx.(store.Tx).Job().Update(ctx, jobModel); der != nil {
			return der
		}

		// if we updated to MINED, we need to update the children and sibling jobs to NEVER_MINED
		if parentJobUUID != "" && nextStatus == entities.StatusMined {
			der := uc.updateChildrenUseCase.
				WithDBTransaction(tx.(store.Tx)).
				Execute(ctx, jobModel.UUID, parentJobUUID, entities.StatusNeverMined, tenants)
			if der != nil {
				return der
			}
		}

		// Metrics observe request latency over job status changes
		if nextStatus != "" && len(jobModel.Logs) > 2 {
			uc.addMetrics(jobModel.Logs[len(jobModel.Logs)-2], jobModel.Logs[len(jobModel.Logs)-1], jobModel.ChainUUID)
		}

		return nil
	})

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	if (nextStatus == entities.StatusMined || nextStatus == entities.StatusStored) && jobModel.NextJobUUID != "" {
		err = uc.startNextJobUseCase.Execute(ctx, jobModel.UUID, tenants)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
		}
	}

	// Refresh job from DB state
	jobModel, err = uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	logger.WithField("status", nextStatus).Info("job updated successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}

func updateJobModel(jobModel *models.Job, nextJob *entities.Job) {
	if len(nextJob.Labels) > 0 {
		jobModel.Labels = nextJob.Labels
	}
	if nextJob.InternalData != nil {
		jobModel.InternalData = nextJob.InternalData
	}

	lastLogID := -1
	for idx, logModel := range jobModel.Logs {
		// Ignore resending and warning statuses
		if logModel.Status == entities.StatusResending || logModel.Status == entities.StatusWarning {
			continue
		}
		// Ignore fail statuses if they come after a resending
		if logModel.Status == entities.StatusFailed && idx > 1 && jobModel.Logs[idx-1].Status == entities.StatusResending {
			continue
		}

		if logModel.ID > lastLogID {
			jobModel.Status = logModel.Status
			lastLogID = logModel.ID
		}
	}
}

func canUpdateStatus(nextStatus, status entities.JobStatus) bool {
	switch nextStatus {
	case entities.StatusCreated:
		return false
	case entities.StatusStarted:
		return status == entities.StatusCreated
	case entities.StatusPending:
		return status == entities.StatusStarted || status == entities.StatusRecovering
	case entities.StatusResending:
		return status == entities.StatusPending
	case entities.StatusRecovering:
		return status == entities.StatusStarted || status == entities.StatusRecovering || status == entities.StatusPending
	case entities.StatusMined, entities.StatusNeverMined:
		return status == entities.StatusPending
	case entities.StatusStored:
		return status == entities.StatusStarted || status == entities.StatusRecovering
	case entities.StatusFailed:
		return status == entities.StatusStarted || status == entities.StatusRecovering || status == entities.StatusPending
	default: // For warning, they can be added at any time
		return true
	}
}

func (uc *updateJobUseCase) addMetrics(current, previous *models.Log, chainUUID string) {
	baseLabels := []string{
		"chain_uuid", chainUUID,
	}

	d := float64(current.CreatedAt.Sub(previous.CreatedAt).Nanoseconds()) / float64(time.Second)
	switch current.Status {
	case entities.StatusMined:
		uc.metrics.MinedLatencyHistogram().With(append(baseLabels,
			"prev_status", string(previous.Status),
			"status", string(current.Status),
		)...).Observe(d)
	default:
		uc.metrics.JobsLatencyHistogram().With(append(baseLabels,
			"prev_status", string(previous.Status),
			"status", string(current.Status),
		)...).Observe(d)
	}

}
