package jobs

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
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
	logger := log.WithContext(ctx).WithField("filters", filters).WithField("tenants", tenants)
	logger.Debug("searching jobs")

	jobModels, err := uc.db.Job().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchJobsComponent)
	}

	var resp []*entities.Job
	for _, jobModel := range jobModels {
		job := parsers.NewJobEntityFromModels(jobModel)
		// Job.Status is a computed value, so that, we filter after parsing
		if filters.Status == "" || job.Status == filters.Status {
			resp = append(resp, job)
		}
	}

	logger.Debug("jobs found successfully")
	return resp, nil
}
