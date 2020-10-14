package client

import (
	"context"

	healthz "github.com/heptiolabs/healthcheck"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type EthereumAccountClient interface {
	CreateAccount(ctx context.Context, request *types.CreateETHAccountRequest) (*types.ETHAccountResponse, error)
}

type KeyManagerClient interface {
	Checker() healthz.Check
	EthereumAccountClient
}
