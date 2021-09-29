package jobs

import (
	"context"

	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store"
)

const getJobComponent = "use-cases.get-job"

// getJobUseCase is a use case to get a job
type getJobUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewGetJobUseCase creates a new GetJobUseCase
func NewGetJobUseCase(db store.DB) usecases.GetJobUseCase {
	return &getJobUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(getJobComponent),
	}
}

// Execute gets a job
func (uc *getJobUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) (*entities.Job, error) {
	ctx = log.WithFields(ctx, log.Field("job", jobUUID))
	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants, true)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getJobComponent)
	}

	uc.logger.WithContext(ctx).Debug("job found successfully")
	return parsers.NewJobEntityFromModels(jobModel), nil
}
