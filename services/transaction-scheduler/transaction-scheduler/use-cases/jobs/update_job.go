package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=update_job.go -destination=mocks/update_job.go -package=mocks

const updateJobComponent = "use-cases.update-job"

type UpdateJobUseCase interface {
	Execute(ctx context.Context, jobEntity *types.Job, newStatus, tenantID string) (*types.Job, error)
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
func (uc *updateJobUseCase) Execute(ctx context.Context, jobEntity *types.Job, newStatus, tenantID string) (*types.Job, error) {
	log.WithContext(ctx).
		WithField("tenant_id", tenantID).
		WithField("job_uuid", jobEntity.UUID).
		Debug("update job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobEntity.UUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		// We are not forced to update the transaction
		if jobEntity.Transaction != nil {
			parsers.UpdateJobModelFromEntities(jobModel, jobEntity)
			if der := tx.(store.Tx).Transaction().Update(ctx, jobModel.Transaction); der != nil {
				return der
			}
			if der := tx.(store.Tx).Job().Update(ctx, jobModel); der != nil {
				return der
			}
		}

		// We are not forced to update the status
		if newStatus != "" {
			jobLogModel := &models.Log{
				JobID:   &jobModel.ID,
				Status:  newStatus,
				Message: "Job updated",
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

	jobModel, err = uc.db.Job().FindOneByUUID(ctx, jobEntity.UUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	log.WithContext(ctx).
		WithField("job_uuid", jobEntity.UUID).
		Info("job updated successfully")

	return parsers.NewJobEntityFromModels(jobModel), nil
}
