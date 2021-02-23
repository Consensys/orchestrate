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

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants, true)
	if err != nil {
		return nil, err
	}

	if entities.IsFinalJobStatus(jobModel.Status) {
		errMessage := "job status is final, cannot be updated"
		logger.WithField("status", jobModel.Status).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage).ExtendComponent(updateJobComponent)
	}

	// We are not forced to update the transaction
	if job.Transaction != nil {
		parsers.UpdateTransactionModelFromEntities(jobModel.Transaction, job.Transaction)
		if err = uc.db.Transaction().Update(ctx, jobModel.Transaction); err != nil {
			return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
		}
	}

	if len(job.Labels) > 0 {
		jobModel.Labels = job.Labels
	}
	if job.InternalData != nil {
		jobModel.InternalData = job.InternalData
	}

	var jobLogModel *models.Log
	// We are not forced to update the status
	if nextStatus != "" && !canUpdateStatus(nextStatus, jobModel.Status) {
		errMessage := "invalid status update for the current job state"
		logger.WithField("status", jobModel.Status).WithField("next_status", nextStatus).Error(errMessage)
		return nil, errors.InvalidStateError(errMessage).ExtendComponent(updateJobComponent)
	} else if nextStatus != "" {
		jobLogModel = &models.Log{
			JobID:   &jobModel.ID,
			Status:  nextStatus,
			Message: logMessage,
		}
	}

	// In case of status update
	err = uc.updateJob(ctx, jobModel, jobLogModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	if (nextStatus == entities.StatusMined || nextStatus == entities.StatusStored) && jobModel.NextJobUUID != "" {
		err = uc.startNextJobUseCase.Execute(ctx, jobModel.UUID, tenants)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
		}
	}

	logger.Info("job updated successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}

func (uc *updateJobUseCase) updateJob(ctx context.Context, jobModel *models.Job, jobLogModel *models.Log) error {
	logger := uc.logger.WithContext(ctx)

	// Does current job belong to a parent/children chains?
	var parentJobUUID string
	if jobModel.InternalData.ParentJobUUID != "" {
		parentJobUUID = jobModel.InternalData.ParentJobUUID
	} else if jobModel.InternalData.RetryInterval != 0 {
		parentJobUUID = jobModel.UUID
	}

	prevLogModel := jobModel.Logs[len(jobModel.Logs)-1]
	err := database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		// We should lock ONLY when there is children jobs
		if parentJobUUID != "" {
			logger.WithField("parent_job", parentJobUUID).Debug("lock parent job row for update")
			if err := tx.(store.Tx).Job().LockOneByUUID(ctx, parentJobUUID); err != nil {
				return err
			}

			// Refresh jobModel after lock to ensure nothing was updated
			refreshedJobModel, err := uc.db.Job().FindOneByUUID(ctx, jobModel.UUID, []string{}, false)
			if err != nil {
				return err
			}

			if refreshedJobModel.UpdatedAt != jobModel.UpdatedAt {
				errMessage := "job status was updated since user request was sent"
				logger.WithField("status", jobModel.Status).Error(errMessage)
				return errors.InvalidStateError(errMessage).ExtendComponent(updateJobComponent)
			}
		}

		if jobLogModel != nil {
			if err := tx.(store.Tx).Log().Insert(ctx, jobLogModel); err != nil {
				return err
			}

			jobModel.Logs = append(jobModel.Logs, jobLogModel)
			if updateNextJobStatus(prevLogModel.Status, jobLogModel.Status) {
				jobModel.Status = jobLogModel.Status
			}
		}

		if err := tx.(store.Tx).Job().Update(ctx, jobModel); err != nil {
			return err
		}

		// if we updated to MINED, we need to update the children and sibling jobs to NEVER_MINED
		if parentJobUUID != "" && jobLogModel != nil && jobLogModel.Status == entities.StatusMined {
			der := uc.updateChildrenUseCase.
				WithDBTransaction(tx.(store.Tx)).
				Execute(ctx, jobModel.UUID, parentJobUUID, entities.StatusNeverMined, []string{})
			if der != nil {
				return der
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Metrics observe request latency over job status changes
	if jobLogModel != nil {
		uc.addMetrics(jobModel.UpdatedAt.Sub(prevLogModel.CreatedAt), prevLogModel.Status, jobLogModel.Status, jobModel.ChainUUID)
	}

	return nil
}

func updateNextJobStatus(prevStatus, nextStatus entities.JobStatus) bool {
	if nextStatus == entities.StatusResending {
		return false
	}
	if nextStatus == entities.StatusWarning {
		return false
	}
	if nextStatus == entities.StatusFailed && prevStatus == entities.StatusResending {
		return false
	}

	return true
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

func (uc *updateJobUseCase) addMetrics(elapseTime time.Duration, previousStatus, nextStatus entities.JobStatus, chainUUID string) {
	if previousStatus == nextStatus {
		return
	}

	baseLabels := []string{
		"chain_uuid", chainUUID,
	}

	switch nextStatus {
	case entities.StatusMined:
		uc.metrics.MinedLatencyHistogram().With(append(baseLabels,
			"prev_status", string(previousStatus),
			"status", string(nextStatus),
		)...).Observe(elapseTime.Seconds())
	default:
		uc.metrics.JobsLatencyHistogram().With(append(baseLabels,
			"prev_status", string(previousStatus),
			"status", string(nextStatus),
		)...).Observe(elapseTime.Seconds())
	}

}
