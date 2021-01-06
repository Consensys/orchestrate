package contracts

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getCatalogComponent = "use-cases.get-catalog"

type getCatalogUseCase struct {
	agent store.RepositoryAgent
}

func NewGetCatalogUseCase(agent store.RepositoryAgent) usecases.GetContractsCatalogUseCase {
	return &getCatalogUseCase{
		agent: agent,
	}
}

// TODO: Modify to get all contracts and then only return necessary fields instead of getting only names
// Execute gets all contract names from DB
func (usecase *getCatalogUseCase) Execute(ctx context.Context) ([]string, error) {
	log.WithContext(ctx).Debug("get catalog starting...")
	names, err := usecase.agent.FindAll(ctx)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getCatalogComponent)
	}

	log.WithContext(ctx).Debug("get catalog executed successfully")
	return names, nil
}
