package types

import "context"

//go:generate mockgen -source=store.go -destination=../mocks/mock_store.go -package=mocks

type ChainStore interface {
	RegisterChain(ctx context.Context, chain *Chain) error

	GetChains(ctx context.Context, filters map[string]string) ([]*Chain, error)
	GetChainsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*Chain, error)
	GetChainByUUID(ctx context.Context, uuid string) (*Chain, error)
	GetChainByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) (*Chain, error)

	UpdateChainByName(ctx context.Context, chain *Chain) error
	UpdateChainByUUID(ctx context.Context, chain *Chain) error

	DeleteChainByUUID(ctx context.Context, uuid string) error
	DeleteChainByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) error
}

type FaucetStore interface {
	RegisterFaucet(ctx context.Context, faucet *Faucet) error

	GetFaucets(ctx context.Context, filters map[string]string) ([]*Faucet, error)
	GetFaucetsByTenant(ctx context.Context, filters map[string]string, tenantID string) ([]*Faucet, error)
	GetFaucetByUUID(ctx context.Context, uuid string) (*Faucet, error)
	GetFaucetByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) (*Faucet, error)

	UpdateFaucetByUUID(ctx context.Context, faucet *Faucet) error

	DeleteFaucetByUUID(ctx context.Context, uuid string) error
	DeleteFaucetByUUIDAndTenant(ctx context.Context, uuid string, tenantID string) error
}

type ChainRegistryStore interface {
	ChainStore
	FaucetStore
}
