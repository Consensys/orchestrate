package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
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
