package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

//go:generate mockgen -source=search_jobs.go -destination=mocks/search_jobs.go -package=mocks

const searchJobsComponent = "use-cases.search-jobs"

type SearchJobsUseCase interface {
	Execute(ctx context.Context, filters map[string]string, tenantID string) ([]*types.JobResponse, error)
}

// searchJobsUseCase is a use case to get a schedule
type searchJobsUseCase struct {
	db store.DB
}

// NewSearchJobsUseCase creates a new SearchJobsUseCase
func NewSearchJobsUseCase(db store.DB) SearchJobsUseCase {
	return &searchJobsUseCase{
		db: db,
	}
}

// Execute gets a schedule
func (uc *searchJobsUseCase) Execute(ctx context.Context, filters map[string]string, tenantID string) ([]*types.JobResponse, error) {
	log.WithContext(ctx).
		WithField("filters", filters).
		Debug("search jobs")

	jobs, err := uc.db.Job().Search(ctx, filters, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchJobsComponent)
	}

	resp := []*types.JobResponse{}
	for _, j := range jobs {
		resp = append(resp, utils.FormatJobResponse(j))
	}

	log.WithContext(ctx).
		WithField("filters", filters).
		Info("jobs found successfully")
	return resp, nil
}
