package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type Client interface {
	GetNodes(ctx context.Context) ([]*types.Node, error)
	GetNodeByID(ctx context.Context, nodeID string) (*types.Node, error)
	GetNodeByTenantAndNodeName(ctx context.Context, tenantID, nodeName string) (*types.Node, error)
	GetNodeByTenantAndNodeID(ctx context.Context, tenantID, nodeID string) (*types.Node, error)
	UpdateBlockPosition(ctx context.Context, nodeID string, blockNumber int64) error
}
