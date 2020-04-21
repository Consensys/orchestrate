package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

//go:generate mockgen -source=data-agents.go -destination=mock/data-agents.go -package=mock

type DataAgents struct {
	Chain     ChainAgent
	Faucet    FaucetAgent
	PrivateTx PrivateTxAgent
}

type ChainAgent interface {
	RegisterChain(ctx context.Context, chain *models.Chain) error

	GetChains(ctx context.Context, filters map[string]string) ([]*models.Chain, error)
	GetChainsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*models.Chain, error)
	GetChainByUUID(ctx context.Context, uuid string) (*models.Chain, error)
	GetChainByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) (*models.Chain, error)

	UpdateChainByName(ctx context.Context, name string, chain *models.Chain) error
	UpdateChainByUUID(ctx context.Context, uuid string, chain *models.Chain) error

	DeleteChainByUUID(ctx context.Context, uuid string) error
	DeleteChainByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) error
}

type FaucetAgent interface {
	RegisterFaucet(ctx context.Context, faucet *models.Faucet) error

	GetFaucets(ctx context.Context, filters map[string]string) ([]*models.Faucet, error)
	GetFaucetsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*models.Faucet, error)
	GetFaucetByUUID(ctx context.Context, uuid string) (*models.Faucet, error)
	GetFaucetByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) (*models.Faucet, error)

	UpdateFaucetByUUID(ctx context.Context, uuid string, faucet *models.Faucet) error

	DeleteFaucetByUUID(ctx context.Context, uuid string) error
	DeleteFaucetByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) error
}

type PrivateTxAgent interface {
	InsertMultiple(ctx context.Context, privateTxManager *[]*models.PrivateTxManagerModel) error
}
