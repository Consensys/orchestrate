package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"
)

//go:generate mockgen -source=update_job.go -destination=mocks/update_job.go -package=mocks

const updateJobComponent = "use-cases.update-job"

// updateJobUseCase is a use case to create a new transaction job
type updateJobUseCase struct {
	db                    store.DB
	updateChildrenUseCase usecases.UpdateChildrenUseCase
	startNextJobUseCase   usecases.StartNextJobUseCase
}

// NewUpdateJobUseCase creates a new UpdateJobUseCase
func NewUpdateJobUseCase(db store.DB, updateChildrenUseCase usecases.UpdateChildrenUseCase, startJobUC usecases.StartNextJobUseCase) usecases.UpdateJobUseCase {
	return &updateJobUseCase{
		db:                    db,
		updateChildrenUseCase: updateChildrenUseCase,
		startNextJobUseCase:   startJobUC,
	}
}

// Execute validates and creates a new transaction job
func (uc *updateJobUseCase) Execute(ctx context.Context, job *entities.Job, nextStatus, logMessage string, tenants []string) (*entities.Job, error) {
	logger := log.WithContext(ctx).WithField("tenants", tenants).WithField("job_uuid", job.UUID)
	logger.Debug("updating job entity")

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
		status := retrievedJob.GetStatus()

		if isFinalStatus(status) {
			errMessage := "job status is final, cannot be updated"
			logger.WithField("status", status).Error(errMessage)
			return errors.InvalidParameterError(errMessage)
		}

		// We are not forced to update the status
		if nextStatus != "" && !canUpdateStatus(nextStatus, status) {
			errMessage := "invalid status update for the current job state"
			logger.WithField("status", status).WithField("next_status", nextStatus).Error(errMessage)
			return errors.InvalidStateError(errMessage)
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
			if nextStatus == utils.StatusMined {
				der := uc.updateChildrenUseCase.
					WithDBTransaction(tx.(store.Tx)).
					Execute(ctx, jobModel.UUID, jobModel.InternalData.ParentJobUUID, utils.StatusNeverMined, tenants)
				if der != nil {
					return der
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	if (nextStatus == utils.StatusMined || nextStatus == utils.StatusStored) && retrievedJob.NextJobUUID != "" {
		err = uc.startNextJobUseCase.Execute(ctx, retrievedJob.UUID, tenants)
		if err != nil {
			logger.WithField("next_job_uuid", retrievedJob.NextJobUUID).WithError(err).Error("fail to start next job")
			return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
		}
	}

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	log.WithContext(ctx).WithField("job_uuid", job.UUID).
		WithField("status", nextStatus).
		Info("job updated successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}

func isFinalStatus(status string) bool {
	return status == utils.StatusMined || status == utils.StatusFailed || status == utils.StatusStored || status == utils.StatusNeverMined
}

func updateJobModel(jobModel *models.Job, job *entities.Job) {
	if len(job.Labels) > 0 {
		jobModel.Labels = job.Labels
	}
	if job.InternalData != nil {
		jobModel.InternalData = job.InternalData
	}
}

func canUpdateStatus(nextStatus, status string) bool {
	switch nextStatus {
	case utils.StatusCreated:
		return false
	case utils.StatusStarted:
		return status == utils.StatusCreated
	case utils.StatusPending:
		return status == utils.StatusStarted || status == utils.StatusRecovering
	case utils.StatusResending:
		return status == utils.StatusPending
	case utils.StatusRecovering, utils.StatusMined, utils.StatusStored, utils.StatusNeverMined:
		return status == utils.StatusPending
	case utils.StatusFailed:
		return status == utils.StatusStarted || status == utils.StatusRecovering || status == utils.StatusPending
	default: // For warning, they can be added at any time
		return true
	}
}
