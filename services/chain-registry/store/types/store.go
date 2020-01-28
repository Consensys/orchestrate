package types

import "context"

type ChainRegistryStore interface {
	RegisterChain(ctx context.Context, chain *Chain) error
	GetChains(ctx context.Context, filters map[string]string) ([]*Chain, error)
	GetChainsByTenantID(ctx context.Context, tenantID string, filters map[string]string) ([]*Chain, error)
	GetChainByTenantIDAndName(ctx context.Context, tenantID string, name string) (*Chain, error)
	GetChainByTenantIDAndUUID(ctx context.Context, tenantID string, uuid string) (*Chain, error)
	GetChainByUUID(ctx context.Context, uuid string) (*Chain, error)
	UpdateChainByName(ctx context.Context, chain *Chain) error
	UpdateBlockPositionByName(ctx context.Context, name, tenantID string, blockPosition int64) error
	UpdateChainByUUID(ctx context.Context, chain *Chain) error
	UpdateBlockPositionByUUID(ctx context.Context, uuid string, blockPosition int64) error
	DeleteChainByName(ctx context.Context, chain *Chain) error
	DeleteChainByUUID(ctx context.Context, uuid string) error
}
