package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

//go:generate mockgen -source=data-agents.go -destination=mock/data-agents.go -package=mock

type DataAgents struct {
	Chain  ChainAgent
	Faucet FaucetAgent
}

type ChainAgent interface {
	RegisterChain(ctx context.Context, chain *models.Chain) error

	GetChains(ctx context.Context, tenants []string, filters map[string]string) ([]*models.Chain, error)
	GetChain(ctx context.Context, uuid string, tenants []string) (*models.Chain, error)
	UpdateChain(ctx context.Context, uuid string, tenants []string, chain *models.Chain) error
	UpdateChainByName(ctx context.Context, name string, tenants []string, chain *models.Chain) error
	DeleteChain(ctx context.Context, uuid string, tenants []string) error
}

type FaucetAgent interface {
	RegisterFaucet(ctx context.Context, faucet *models.Faucet) error
	GetFaucets(ctx context.Context, tenants []string, filters map[string]string) ([]*models.Faucet, error)
	GetFaucet(ctx context.Context, uuid string, tenants []string) (*models.Faucet, error)
	UpdateFaucet(ctx context.Context, uuid string, tenants []string, faucet *models.Faucet) error
	DeleteFaucet(ctx context.Context, uuid string, tenants []string) error
}
