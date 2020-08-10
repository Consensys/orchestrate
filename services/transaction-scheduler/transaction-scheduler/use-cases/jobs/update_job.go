package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=update_job.go -destination=mocks/update_job.go -package=mocks

const updateJobComponent = "use-cases.update-job"

type UpdateJobUseCase interface {
	Execute(ctx context.Context, jobEntity *types.Job, nextStatus, logMessage string, tenants []string) (*types.Job, error)
}

// updateJobUseCase is a use case to create a new transaction job
type updateJobUseCase struct {
	db store.DB
}

// NewUpdateJobUseCase creates a new UpdateJobUseCase
func NewUpdateJobUseCase(db store.DB) UpdateJobUseCase {
	return &updateJobUseCase{
		db: db,
	}
}

// Execute validates and creates a new transaction job
func (uc *updateJobUseCase) Execute(ctx context.Context, job *types.Job, nextStatus, logMessage string, tenants []string) (*types.Job, error) {
	logger := log.WithContext(ctx).WithField("tenants", tenants).WithField("job_uuid", job.UUID)
	logger.Debug("updating job entity")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	retrievedJob := parsers.NewJobEntityFromModels(jobModel)
	status := retrievedJob.GetStatus()
	if status == utils.StatusMined || status == utils.StatusFailed {
		errMessage := "job cannot be updated in the current state"
		logger.WithField("status", status).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage).ExtendComponent(updateJobComponent)
	}

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
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
			if !canUpdateStatus(nextStatus, status) {
				errMessage := "invalid status update for the current job state"
				logger.WithField("status", status).WithField("next_status", nextStatus).Error(errMessage)
				return errors.InvalidStateError(errMessage)
			}

			jobLogModel := &models.Log{
				JobID:   &jobModel.ID,
				Status:  nextStatus,
				Message: logMessage,
			}
			if der := tx.(store.Tx).Log().Insert(ctx, jobLogModel); der != nil {
				return der
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	jobModel, err = uc.db.Job().FindOneByUUID(ctx, job.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	log.WithContext(ctx).WithField("job_uuid", job.UUID).Info("job updated successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}

func updateJobModel(jobModel *models.Job, job *types.Job) {
	if len(job.Labels) > 0 {
		jobModel.Labels = job.Labels
	}
	if job.Annotations != nil {
		jobModel.Annotations = job.Annotations
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
	case utils.StatusRecovering, utils.StatusMined:
		return status == utils.StatusPending
	case utils.StatusFailed:
		return status == utils.StatusStarted || status == utils.StatusRecovering || status == utils.StatusPending
	default: // For warning, they can be added at any time
		return true
	}
}
