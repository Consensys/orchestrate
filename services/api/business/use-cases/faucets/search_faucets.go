package faucets

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
)

const searchFaucetsComponent = "use-cases.search-faucets"

// searchFaucetsUseCase is a use case to search faucets
type searchFaucetsUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewSearchFaucets creates a new SearchFaucetsUseCase
func NewSearchFaucets(db store.DB) usecases.SearchFaucetsUseCase {
	return &searchFaucetsUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(searchFaucetsComponent),
	}
}

// Execute search faucets
func (uc *searchFaucetsUseCase) Execute(ctx context.Context, filters *entities.FaucetFilters, userInfo *multitenancy.UserInfo) ([]*entities.Faucet, error) {
	faucetModels, err := uc.db.Faucet().Search(ctx, filters, userInfo.AllowedTenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchFaucetsComponent)
	}

	var faucets []*entities.Faucet
	for _, faucetModel := range faucetModels {
		faucets = append(faucets, parsers.NewFaucetFromModel(faucetModel))
	}

	uc.logger.Debug("faucets found successfully")
	return faucets, nil
}
