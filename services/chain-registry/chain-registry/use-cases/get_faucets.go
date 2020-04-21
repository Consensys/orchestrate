package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type GetFaucets interface {
	Execute(ctx context.Context, tenantID string, filters map[string]string) ([]*models.Faucet, error)
}

// RegisterContract is a use case to register a new contract
type getFaucets struct {
	faucetAgent store.FaucetAgent
}

// NewGetCatalog creates a new GetCatalog
func NewGetFaucets(faucetAgent store.FaucetAgent) GetFaucets {
	return &getFaucets{
		faucetAgent: faucetAgent,
	}
}

func (uc *getFaucets) Execute(ctx context.Context, tenantID string, filters map[string]string) ([]*models.Faucet, error) {
	var faucets []*models.Faucet
	var err error

	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		faucets, err = uc.faucetAgent.GetFaucets(ctx, filters)
	} else {
		faucets, err = uc.faucetAgent.GetFaucetsByTenant(ctx, filters, tenantID)
	}

	return faucets, err
}
