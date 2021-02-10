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
	ctx = log.WithFields(ctx, log.Field("job", job.UUID))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("updating job")

	var retrievedJob *entities.Job
	err := database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		der := tx.(store.Tx).Job().LockOneByUUID(ctx, job.UUID)
		if der != nil {
			return der
		}

		jobModel, der := tx.(store.Tx).Job().FindOneByUUID(ctx, job.UUID, tenants)
		if der != nil {
			return der
		}

		retrievedJob = parsers.NewJobEntityFromModels(jobModel)
		status := retrievedJob.Status

		if isFinalStatus(status) {
			errMessage := "job status is final, cannot be updated"
			logger.WithField("status", status).Error(errMessage)
			return errors.InvalidParameterError(errMessage).ExtendComponent(updateJobComponent)
		}

		// We are not forced to update the status
		if nextStatus != "" && !canUpdateStatus(nextStatus, status) {
			errMessage := "invalid status update for the current job state"
			logger.WithField("status", status).WithField("next_status", nextStatus).Error(errMessage)
			return errors.InvalidStateError(errMessage).ExtendComponent(updateJobComponent)
		}

		// We are not forced to update the transaction
		if job.Transaction != nil {
			parsers.UpdateTransactionModelFromEntities(jobModel.Transaction, job.Transaction)
			if der := tx.(store.Tx).Transaction().Update(ctx, jobModel.Transaction); der != nil {
				return der
			}
		}

		updateJobModel(jobModel, job)
		if der := tx.(store.Tx).Job().Update(ctx, jobModel); der != nil {
			return der
		}

		// We are not forced to update the status
		if nextStatus != "" {
			jobLogModel := &models.Log{
				JobID:   &jobModel.ID,
				Status:  nextStatus,
				Message: logMessage,
			}
			if der := tx.(store.Tx).Log().Insert(ctx, jobLogModel); der != nil {
				return der
			}

			// if we updated to MINED, we need to update the children and sibling jobs to NEVER_MINED
			if nextStatus == entities.StatusMined {
				der := uc.updateChildrenUseCase.
					WithDBTransaction(tx.(store.Tx)).
					Execute(ctx, jobModel.UUID, jobModel.InternalData.ParentJobUUID, entities.StatusNeverMined, tenants)
				if der != nil {
					return der
				}
			}

			// Metrics observe request latency
			uc.addMetrics(jobLogModel, jobModel.Logs[len(jobModel.Logs)-1], jobModel.ChainUUID)
		}

		return nil
	})

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	if (nextStatus == entities.StatusMined || nextStatus == entities.StatusStored) && retrievedJob.NextJobUUID != "" {
		err = uc.startNextJobUseCase.Execute(ctx, retrievedJob.UUID, tenants)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
		}
	}

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	logger.WithField("status", nextStatus).Info("job updated successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}

func isFinalStatus(status entities.JobStatus) bool {
	return status == entities.StatusMined ||
		status == entities.StatusFailed ||
		status == entities.StatusStored ||
		status == entities.StatusNeverMined
}

func updateJobModel(jobModel *models.Job, job *entities.Job) {
	if len(job.Labels) > 0 {
		jobModel.Labels = job.Labels
	}
	if job.InternalData != nil {
		jobModel.InternalData = job.InternalData
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
