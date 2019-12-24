package client

import "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

type Client interface {
	GetNodeByID(id string) (*types.Node, error)
	GetNodes() ([]*types.Node, error)
}
