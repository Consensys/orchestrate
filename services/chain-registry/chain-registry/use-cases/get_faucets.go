package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type GetFaucets interface {
	Execute(ctx context.Context, tenants []string, filters map[string]string) ([]*models.Faucet, error)
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

func (uc *getFaucets) Execute(ctx context.Context, tenants []string, filters map[string]string) ([]*models.Faucet, error) {
	return uc.faucetAgent.GetFaucets(ctx, tenants, filters)
}
