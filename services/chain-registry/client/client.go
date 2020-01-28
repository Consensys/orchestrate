package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type Client interface {
	GetChains(ctx context.Context) ([]*types.Chain, error)
	GetChainByUUID(ctx context.Context, chainUUID string) (*types.Chain, error)
	GetChainByTenantAndName(ctx context.Context, tenantID, chainName string) (*types.Chain, error)
	GetChainByTenantAndUUID(ctx context.Context, tenantID, chainUUID string) (*types.Chain, error)
	UpdateBlockPosition(ctx context.Context, chainUUID string, blockNumber int64) error
}
