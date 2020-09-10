package faucets

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type GetFaucet interface {
	Execute(ctx context.Context, uuid string, tenants []string) (*models.Faucet, error)
}

// RegisterContract is a use case to register a new contract
type getFaucet struct {
	faucetAgent store.FaucetAgent
}

// NewGetCatalog creates a new GetCatalog
func NewGetFaucet(faucetAgent store.FaucetAgent) GetFaucet {
	return &getFaucet{
		faucetAgent: faucetAgent,
	}
}

func (uc *getFaucet) Execute(ctx context.Context, uuid string, tenants []string) (*models.Faucet, error) {
	return uc.faucetAgent.GetFaucet(ctx, uuid, tenants)
}
