package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

// IMPORTANT: Mock is created in a separated folder because of cycle deps
// https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/pegasyseng/orchestrate/296

//go:generate mockgen -source=create_job.go -destination=mocks/create_job.go -package=mocks

const createJobComponent = "use-cases.create-job"

// createJobUseCase is a use case to create a new transaction job
type createJobUseCase struct {
	validator validators.TransactionValidator
	db        store.DB
}

// NewCreateJobUseCase creates a new CreateJobUseCase
func NewCreateJobUseCase(db store.DB, validator validators.TransactionValidator) usecases.CreateJobUseCase {
	return &createJobUseCase{
		validator: validator,
		db:        db,
	}
}

func (uc createJobUseCase) WithDBTransaction(dbtx store.Tx) usecases.CreateJobUseCase {
	uc.db = dbtx
	return &uc
}

// Execute validates and creates a new transaction job
func (uc *createJobUseCase) Execute(ctx context.Context, job *entities.Job, tenants []string) (*entities.Job, error) {
	logger := log.WithContext(ctx).
		WithField("chain_uuid", job.ChainUUID).
		WithField("schedule_id", job.ScheduleUUID).
		WithField("tenants", tenants)
	logger.Debug("creating new job")

	chainID, err := uc.validator.ValidateChainExists(ctx, job.ChainUUID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}
	job.InternalData.ChainID = chainID

	schedule, err := uc.db.Schedule().FindOneByUUID(ctx, job.ScheduleUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	jobModel := parsers.NewJobModelFromEntities(job, &schedule.ID)
	jobModel.Logs = append(jobModel.Logs, &models.Log{
		Status: utils.StatusCreated,
	})

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		// If it's a child job, only create it if parent status is PENDING
		if jobModel.InternalData.ParentJobUUID != "" {
			if der := tx.(store.Tx).Job().LockOneByUUID(ctx, jobModel.InternalData.ParentJobUUID); der != nil {
				return der
			}

			parentJobModel, der := tx.(store.Tx).Job().FindOneByUUID(ctx, jobModel.InternalData.ParentJobUUID, tenants)
			if der != nil {
				return der
			}

			parentStatus := parsers.NewJobEntityFromModels(parentJobModel).GetStatus()
			if parentStatus != utils.StatusPending {
				errMessage := "cannot create a child job in a finalized schedule"
				logger.
					WithField("parent_job_uuid", jobModel.InternalData.ParentJobUUID).
					WithField("parent_status", parentStatus).
					Error(errMessage)
				return errors.InvalidStateError(errMessage)
			}
		}

		if der := tx.(store.Tx).Transaction().Insert(ctx, jobModel.Transaction); der != nil {
			return der
		}

		if der := tx.(store.Tx).Job().Insert(ctx, jobModel); der != nil {
			return der
		}

		jobModel.Logs[0].JobID = &jobModel.ID
		if der := tx.(store.Tx).Log().Insert(ctx, jobModel.Logs[0]); der != nil {
			return der
		}

		return nil
	})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	log.WithContext(ctx).WithField("job_uuid", jobModel.UUID).Info("job created successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}
