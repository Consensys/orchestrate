package client

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	healthz "github.com/heptiolabs/healthcheck"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type ChainClient interface {
	GetChains(ctx context.Context) ([]*models.Chain, error)
	GetChainByName(ctx context.Context, chainName string) (*models.Chain, error)
	GetChainByUUID(ctx context.Context, chainUUID string) (*models.Chain, error)
	DeleteChainByUUID(ctx context.Context, chainUUID string) error
	RegisterChain(ctx context.Context, chain *models.Chain) (*models.Chain, error)
	UpdateBlockPosition(ctx context.Context, chainUUID string, blockNumber uint64) error
	UpdateChainByUUID(ctx context.Context, chainUUID string, chain *models.Chain) error
}

type FaucetClient interface {
	RegisterFaucet(ctx context.Context, faucet *models.Faucet) (*models.Faucet, error)
	UpdateFaucetByUUID(ctx context.Context, uuid string, faucet *models.Faucet) (*models.Faucet, error)
	GetFaucetByUUID(ctx context.Context, faucetUUID string) (*models.Faucet, error)
	DeleteFaucetByUUID(ctx context.Context, faucetUUID string) error
	GetFaucetsByChainRule(ctx context.Context, chainRule string) ([]*models.Faucet, error)
	GetFaucetCandidate(ctx context.Context, sender ethcommon.Address, chainUUID string) (*types.Faucet, error)
}

type ChainRegistryClient interface {
	Checker() healthz.Check
	ChainClient
	FaucetClient
}
