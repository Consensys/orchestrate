package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=create_job.go -destination=mocks/create_job.go -package=mocks

const createJobComponent = "use-cases.create-job"

type CreateJobUseCase interface {
	Execute(ctx context.Context, job *entities.Job, tenantID string) (*entities.Job, error)
}

// createJobUseCase is a use case to create a new transaction job
type createJobUseCase struct {
	db store.DB
}

// NewCreateJobUseCase creates a new CreateJobUseCase
func NewCreateJobUseCase(db store.DB) CreateJobUseCase {
	return &createJobUseCase{
		db: db,
	}
}

// Execute validates and creates a new transaction job
func (uc *createJobUseCase) Execute(ctx context.Context, job *entities.Job, tenantID string) (*entities.Job, error) {
	log.WithContext(ctx).
		WithField("schedule_id", job.ScheduleUUID).
		WithField("tenant_id", tenantID).
		Debug("creating new job")

	schedule, err := uc.db.Schedule().FindOneByUUID(ctx, job.ScheduleUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	jobModel := parsers.NewJobModelFromEntities(job, &schedule.ID)
	jobModel.Logs = append(jobModel.Logs, &models.Log{
		Status:  entities.JobStatusCreated,
		Message: "Job created",
	})

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
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
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	log.WithContext(ctx).
		WithField("job_uuid", jobModel.UUID).
		Info("job created successfully")

	return parsers.NewJobEntityFromModels(jobModel), nil
}
