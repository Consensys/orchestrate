package jobs

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

const searchJobsComponent = "use-cases.search-jobs"

// searchJobsUseCase is a use case to search jobs
type searchJobsUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewSearchJobsUseCase creates a new SearchJobsUseCase
func NewSearchJobsUseCase(db store.DB) usecases.SearchJobsUseCase {
	return &searchJobsUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(searchJobsComponent),
	}
}

// Execute search jobs
func (uc *searchJobsUseCase) Execute(ctx context.Context, filters *entities.JobFilters, tenants []string) ([]*entities.Job, error) {
	jobModels, err := uc.db.Job().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchJobsComponent)
	}

	var resp []*entities.Job
	for _, jobModel := range jobModels {
		job := parsers.NewJobEntityFromModels(jobModel)
		resp = append(resp, job)
	}

	uc.logger.Debug("jobs found successfully")
	return resp, nil
}
