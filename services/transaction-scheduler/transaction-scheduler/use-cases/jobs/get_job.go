package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=get_job.go -destination=mocks/get_job.go -package=mocks

const getJobComponent = "use-cases.get-job"

type GetJobUseCase interface {
	Execute(ctx context.Context, jobUUID string, tenants []string) (*types.Job, error)
}

// getJobUseCase is a use case to get a job
type getJobUseCase struct {
	db store.DB
}

// NewGetJobUseCase creates a new GetJobUseCase
func NewGetJobUseCase(db store.DB) GetJobUseCase {
	return &getJobUseCase{
		db: db,
	}
}

// Execute gets a job
func (uc *getJobUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) (*types.Job, error) {
	log.WithContext(ctx).
		WithField("job_uuid", jobUUID).
		WithField("tenants", tenants).
		Debug("getting job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getJobComponent)
	}

	log.WithContext(ctx).
		WithField("job_uuid", jobUUID).
		WithField("tenants", tenants).
		Info("job found successfully")

	return parsers.NewJobEntityFromModels(jobModel), nil
}
