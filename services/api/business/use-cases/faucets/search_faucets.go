package faucets

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const searchFaucetsComponent = "use-cases.search-faucets"

// searchFaucetsUseCase is a use case to search faucets
type searchFaucetsUseCase struct {
	db store.DB
}

// NewSearchFaucets creates a new SearchFaucetsUseCase
func NewSearchFaucets(db store.DB) usecases.SearchFaucetsUseCase {
	return &searchFaucetsUseCase{
		db: db,
	}
}

// Execute search faucets
func (uc *searchFaucetsUseCase) Execute(ctx context.Context, filters *entities.FaucetFilters, tenants []string) ([]*entities.Faucet, error) {
	logger := log.WithContext(ctx).WithField("filters", filters).WithField("tenants", tenants)
	logger.Debug("searching faucets")

	faucetModels, err := uc.db.Faucet().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchFaucetsComponent)
	}

	var faucets []*entities.Faucet
	for _, faucetModel := range faucetModels {
		faucets = append(faucets, parsers.NewFaucetFromModel(faucetModel))
	}

	logger.Debug("faucets found successfully")
	return faucets, nil
}
