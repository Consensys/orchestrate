package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const createJobComponent = "use-cases.create-job"

// createJobUseCase is a use case to create a new transaction job
type createJobUseCase struct {
	db         store.DB
	getChainUC usecases.GetChainUseCase
	logger     *log.Logger
}

// NewCreateJobUseCase creates a new CreateJobUseCase
func NewCreateJobUseCase(db store.DB, getChainUC usecases.GetChainUseCase) usecases.CreateJobUseCase {
	return &createJobUseCase{
		db:         db,
		getChainUC: getChainUC,
		logger:     log.NewLogger().SetComponent(createJobComponent),
	}
}

func (uc createJobUseCase) WithDBTransaction(dbtx store.Tx) usecases.CreateJobUseCase {
	uc.db = dbtx
	return &uc
}

// Execute validates and creates a new transaction job
func (uc *createJobUseCase) Execute(ctx context.Context, job *entities.Job, tenants []string) (*entities.Job, error) {
	ctx = log.WithFields(ctx, log.Field("chain", job.ChainUUID), log.Field("schedule", job.ScheduleUUID))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new job")

	chainID, err := uc.getChainID(ctx, job.ChainUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}
	job.InternalData.ChainID = chainID

	if job.Transaction.From != "" {
		err = uc.validateAccountExists(ctx, job.Transaction.From, tenants)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(createJobComponent)
		}
	}

	schedule, err := uc.db.Schedule().FindOneByUUID(ctx, job.ScheduleUUID, tenants)
	if errors.IsNotFoundError(err) {
		return nil, errors.InvalidParameterError("schedule does not exist").ExtendComponent(createJobComponent)
	}
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	jobModel := parsers.NewJobModelFromEntities(job, &schedule.ID)
	jobModel.Logs = append(jobModel.Logs, &models.Log{
		Status: utils.StatusCreated,
	})
	jobModel.Schedule = schedule

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

			parentStatus := parsers.NewJobEntityFromModels(parentJobModel).Status
			if parentStatus != utils.StatusPending {
				errMessage := "cannot create a child job in a finalized schedule"
				logger.WithField("parent_job", jobModel.InternalData.ParentJobUUID).
					WithField("parent_status", parentStatus).Error(errMessage)
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

	logger.WithField("job", jobModel.UUID).Info("job created successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}

func (uc *createJobUseCase) validateAccountExists(ctx context.Context, address string, tenants []string) error {
	_, err := uc.db.Account().FindOneByAddress(ctx, address, tenants)
	if errors.IsNotFoundError(err) {
		return errors.InvalidParameterError("failed to get account")
	}
	if err != nil {
		return err
	}

	return nil
}

func (uc *createJobUseCase) getChainID(ctx context.Context, chainUUID string, tenants []string) (string, error) {
	chain, err := uc.getChainUC.Execute(ctx, chainUUID, tenants)
	if errors.IsNotFoundError(err) {
		return "", errors.InvalidParameterError("failed to get chain")
	}
	if err != nil {
		return "", errors.FromError(err)
	}

	return chain.ChainID, nil
}
