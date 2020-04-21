package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type GetChains interface {
	Execute(ctx context.Context, tenantID string, filters map[string]string) ([]*models.Chain, error)
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

func (uc *getChains) Execute(ctx context.Context, tenantID string, filters map[string]string) ([]*models.Chain, error) {
	var chains []*models.Chain
	var err error

	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		chains, err = uc.chainAgent.GetChains(ctx, filters)
	} else {
		chains, err = uc.chainAgent.GetChainsByTenant(ctx, filters, tenantID)
	}

	return chains, err
}
