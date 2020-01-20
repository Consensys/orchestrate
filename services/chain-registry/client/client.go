package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type Client interface {
	GetNodeByID(ctx context.Context, nodeID string) (*types.Node, error)
	GetNodes(ctx context.Context) ([]*types.Node, error)
	UpdateBlockPosition(ctx context.Context, nodeID string, blockNumber int64) error
}
