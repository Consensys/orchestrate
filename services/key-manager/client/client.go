package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager"

	healthz "github.com/heptiolabs/healthcheck"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type EthereumAccountClient interface {
	ETHCreateAccount(ctx context.Context, request *types.CreateETHAccountRequest) (*types.ETHAccountResponse, error)
	ETHImportAccount(ctx context.Context, request *types.ImportETHAccountRequest) (*types.ETHAccountResponse, error)
	ETHSign(ctx context.Context, address string, request *keymanager.PayloadRequest) (string, error)
	ETHSignTransaction(ctx context.Context, address string, request *types.SignETHTransactionRequest) (string, error)
	ETHSignTesseraTransaction(ctx context.Context, address string, request *types.SignTesseraTransactionRequest) (string, error)
}

type KeyManagerClient interface {
	Checker() healthz.Check
	EthereumAccountClient
}
