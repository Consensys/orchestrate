package chains

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

const searchChainsComponent = "use-cases.search-chains"

// searchChainsUseCase is a use case to search chains
type searchChainsUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewSearchChainsUseCase creates a new SearchChainsUseCase
func NewSearchChainsUseCase(db store.DB) usecases.SearchChainsUseCase {
	return &searchChainsUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(searchChainsComponent),
	}
}

// Execute search faucets
func (uc *searchChainsUseCase) Execute(ctx context.Context, filters *entities.ChainFilters, tenants []string) ([]*entities.Chain, error) {
	logger := uc.logger.WithContext(ctx)

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
