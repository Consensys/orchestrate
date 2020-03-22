package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type ChainClient interface {
	GetChains(ctx context.Context) ([]*types.Chain, error)
	GetChainByName(ctx context.Context, chainName string) (*types.Chain, error)
	GetChainByUUID(ctx context.Context, chainUUID string) (*types.Chain, error)
	UpdateBlockPosition(ctx context.Context, chainUUID string, blockNumber uint64) error
}

type FaucetClient interface {
	GetFaucetsByChainRule(ctx context.Context, chainRule string) ([]*types.Faucet, error)
}

type ChainRegistryClient interface {
	ChainClient
	FaucetClient
}
