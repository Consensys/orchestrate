package client

import (
	"context"

	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
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

type ChainRegistryClient interface {
	Checker() healthz.Check
	ChainClient
}
