package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=get_job.go -destination=mocks/get_job.go -package=mocks

const getJobComponent = "use-cases.get-job"

// getJobUseCase is a use case to get a job
type getJobUseCase struct {
	db store.DB
}

// NewGetJobUseCase creates a new GetJobUseCase
func NewGetJobUseCase(db store.DB) usecases.GetJobUseCase {
	return &getJobUseCase{
		db: db,
	}
}

// Execute gets a job
func (uc *getJobUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) (*entities.Job, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", jobUUID).WithField("tenants", tenants)
	logger.Debug("getting job")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getJobComponent)
	}

	logger.Debug("job found successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}
