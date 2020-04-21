package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type GetChain interface {
	Execute(ctx context.Context, uuid, tenantID string) (*models.Chain, error)
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

func (uc *getChain) Execute(ctx context.Context, uuid, tenantID string) (*models.Chain, error) {
	var chain *models.Chain
	var err error
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		chain, err = uc.chainAgent.GetChainByUUID(ctx, uuid)
	} else {
		chain, err = uc.chainAgent.GetChainByUUIDAndTenant(ctx, uuid, tenantID)
	}

	return chain, err
}
