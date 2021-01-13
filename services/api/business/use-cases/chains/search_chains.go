package chains

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const searchChainsComponent = "use-cases.search-chains"

// searchChainsUseCase is a use case to search chains
type searchChainsUseCase struct {
	db store.DB
}

// NewSearchChainsUseCase creates a new SearchChainsUseCase
func NewSearchChainsUseCase(db store.DB) usecases.SearchChainsUseCase {
	return &searchChainsUseCase{
		db: db,
	}
}

// Execute search faucets
func (uc *searchChainsUseCase) Execute(ctx context.Context, filters *entities.ChainFilters, tenants []string) ([]*entities.Chain, error) {
	logger := log.WithContext(ctx).WithField("filters", filters).WithField("tenants", tenants)
	logger.Debug("searching chains")

	chainModels, err := uc.db.Chain().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchChainsComponent)
	}

	var chains []*entities.Chain
	for _, chainModel := range chainModels {
		chains = append(chains, parsers.NewChainFromModel(chainModel))
	}

	logger.Debug("chains found successfully")
	return chains, nil
}
