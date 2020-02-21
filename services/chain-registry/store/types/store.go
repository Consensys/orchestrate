package types

import "context"

type ChainRegistryStore interface {
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
