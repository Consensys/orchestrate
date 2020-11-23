package chains

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

type GetChain interface {
	Execute(ctx context.Context, uuid string, tenants []string) (*models.Chain, error)
}

// RegisterContract is a use case to register a new contract
type getChain struct {
	chainAgent store.ChainAgent
}

// NewGetCatalog creates a new GetCatalog
func NewGetChain(chainAgent store.ChainAgent) GetChain {
	return &getChain{
		chainAgent: chainAgent,
	}
}

func (uc *getChain) Execute(ctx context.Context, uuid string, tenants []string) (*models.Chain, error) {
	return uc.chainAgent.GetChain(ctx, uuid, tenants)
}
