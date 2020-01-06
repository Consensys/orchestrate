package client

import "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

type Client interface {
	GetNodeByID(nodeID string) (*types.Node, error)
	GetNodes() ([]*types.Node, error)
	UpdateBlockPosition(nodeID string, blockNumber int64) error
}
