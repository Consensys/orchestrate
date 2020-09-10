package chains

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type GetChains interface {
	Execute(ctx context.Context, tenant []string, filters map[string]string) ([]*models.Chain, error)
}

// RegisterContract is a use case to register a new contract
type getChains struct {
	chainAgent store.ChainAgent
}

// NewGetCatalog creates a new GetCatalog
func NewGetChains(chainAgent store.ChainAgent) GetChains {
	return &getChains{
		chainAgent: chainAgent,
	}
}

func (uc *getChains) Execute(ctx context.Context, tenants []string, filters map[string]string) ([]*models.Chain, error) {
	return uc.chainAgent.GetChains(ctx, tenants, filters)
}
