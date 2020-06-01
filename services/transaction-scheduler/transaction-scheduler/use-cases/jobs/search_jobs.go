package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=search_jobs.go -destination=mocks/search_jobs.go -package=mocks

const searchJobsComponent = "use-cases.search-jobs"

type SearchJobsUseCase interface {
	Execute(ctx context.Context, filters *entities.JobFilters, tenantID string) ([]*entities.Job, error)
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
func (uc *searchJobsUseCase) Execute(ctx context.Context, filters *entities.JobFilters, tenantID string) ([]*entities.Job, error) {
	log.WithContext(ctx).
		WithField("filters", filters).
		Debug("search jobs")

	txHashesFilter := []string{}
	for _, hash := range filters.TxHashes {
		txHashesFilter = append(txHashesFilter, hash.String())
	}

	jobModels, err := uc.db.Job().Search(ctx, tenantID, txHashesFilter)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchJobsComponent)
	}

	resp := []*entities.Job{}
	for _, jb := range jobModels {
		resp = append(resp, parsers.NewJobEntityFromModels(jb))
	}

	log.WithContext(ctx).
		WithField("filters", filters).
		Info("jobs found successfully")

	return resp, nil
}
