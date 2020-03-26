package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
)

const getCatalogComponent = component + ".get-catalog"

type GetCatalogUseCase interface {
	Execute(ctx context.Context) ([]string, error)
}

// GetCatalog is a use case to get all contract names
type GetCatalog struct {
	repositoryDataAgent store.RepositoryDataAgent
}

// NewGetCatalog creates a new GetCatalog
func NewGetCatalog(repositoryDataAgent store.RepositoryDataAgent) *GetCatalog {
	return &GetCatalog{
		repositoryDataAgent: repositoryDataAgent,
	}
}

// TODO: Modify to get all contracts and then only return necessary fields instead of getting only names
// Execute gets all contract names from DB
func (usecase *GetCatalog) Execute(ctx context.Context) ([]string, error) {
	names, err := usecase.repositoryDataAgent.FindAll(ctx)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getCatalogComponent)
	}

	return names, nil
}
