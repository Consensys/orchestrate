package jobs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=search_jobs.go -destination=mocks/search_jobs.go -package=mocks

const searchJobsComponent = "use-cases.search-jobs"

type SearchJobsUseCase interface {
	Execute(ctx context.Context, filters *entities.JobFilters, tenants []string) ([]*types.Job, error)
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
func (uc *searchJobsUseCase) Execute(ctx context.Context, filters *entities.JobFilters, tenants []string) ([]*types.Job, error) {
	log.WithContext(ctx).
		WithField("filters", filters).
		Debug("search jobs")

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(searchJobsComponent)
	}

	jobModels, err := uc.db.Job().Search(ctx, filters.TxHashes, filters.ChainUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchJobsComponent)
	}

	var resp []*types.Job
	for _, jb := range jobModels {
		resp = append(resp, parsers.NewJobEntityFromModels(jb))
	}

	log.WithContext(ctx).
		WithField("filters", filters).
		Info("jobs found successfully")

	return resp, nil
}
