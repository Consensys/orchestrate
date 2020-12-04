package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"

	healthz "github.com/heptiolabs/healthcheck"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type EthereumAccountClient interface {
	ETHCreateAccount(ctx context.Context, request *types.CreateETHAccountRequest) (*types.ETHAccountResponse, error)
	ETHImportAccount(ctx context.Context, request *types.ImportETHAccountRequest) (*types.ETHAccountResponse, error)
	ETHSign(ctx context.Context, address string, request *keymanager.PayloadRequest) (string, error)
	ETHSignTypedData(ctx context.Context, address string, request *types.SignTypedDataRequest) (string, error)
	ETHSignTransaction(ctx context.Context, address string, request *types.SignETHTransactionRequest) (string, error)
	ETHSignQuorumPrivateTransaction(ctx context.Context, address string, request *types.SignQuorumPrivateTransactionRequest) (string, error)
	ETHSignEEATransaction(ctx context.Context, address string, request *types.SignEEATransactionRequest) (string, error)
	ETHListAccounts(ctx context.Context, namespace string) ([]string, error)
	ETHListNamespaces(ctx context.Context) ([]string, error)
	ETHGetAccount(ctx context.Context, address, namespace string) (*types.ETHAccountResponse, error)
}

type KeyManagerClient interface {
	Checker() healthz.Check
	EthereumAccountClient
}
