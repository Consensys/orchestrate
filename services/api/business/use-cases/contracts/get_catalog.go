package contracts

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getCatalogComponent = "use-cases.get-catalog"

type getCatalogUseCase struct {
	agent  store.RepositoryAgent
	logger *log.Logger
}

func NewGetCatalogUseCase(agent store.RepositoryAgent) usecases.GetContractsCatalogUseCase {
	return &getCatalogUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(getCatalogComponent),
	}
}

// TODO: Modify to get all contracts and then only return necessary fields instead of getting only names
// Execute gets all contract names from DB
func (uc *getCatalogUseCase) Execute(ctx context.Context) ([]string, error) {
	names, err := uc.agent.FindAll(ctx)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getCatalogComponent)
	}

	uc.logger.WithContext(ctx).Debug("get catalog executed successfully")
	return names, nil
}
