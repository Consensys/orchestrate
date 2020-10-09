package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

const searchJobsComponent = "use-cases.search-jobs"

// searchJobsUseCase is a use case to search jobs
type searchJobsUseCase struct {
	db store.DB
}

// NewSearchJobsUseCase creates a new SearchJobsUseCase
func NewSearchJobsUseCase(db store.DB) usecases.SearchJobsUseCase {
	return &searchJobsUseCase{
		db: db,
	}
}

// Execute search jobs
func (uc *searchJobsUseCase) Execute(ctx context.Context, filters *entities.JobFilters, tenants []string) ([]*entities.Job, error) {
	log.WithContext(ctx).WithField("filters", filters).WithField("tenants", tenants).Debug("searching jobs")

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(searchJobsComponent)
	}

	jobModels, err := uc.db.Job().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchJobsComponent)
	}

	var resp []*entities.Job
	for _, jobModel := range jobModels {
		job := parsers.NewJobEntityFromModels(jobModel)
		// Job.Status is a computed value, so that, we filter after parsing
		if filters.Status == "" || job.GetStatus() == filters.Status {
			resp = append(resp, job)
		}
	}

	// Debug as search jobs is constantly called by tx-listener and tx-sentry
	log.WithContext(ctx).Debug("jobs found successfully")
	return resp, nil
}
