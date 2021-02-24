package client

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/keymanager"
	ethTypes "github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
	zksTypes "github.com/ConsenSys/orchestrate/pkg/types/keymanager/zk-snarks"
	healthz "github.com/heptiolabs/healthcheck"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type EthereumAccountClient interface {
	ETHCreateAccount(ctx context.Context, request *ethTypes.CreateETHAccountRequest) (*ethTypes.ETHAccountResponse, error)
	ETHImportAccount(ctx context.Context, request *ethTypes.ImportETHAccountRequest) (*ethTypes.ETHAccountResponse, error)
	ETHSign(ctx context.Context, address string, request *keymanager.SignPayloadRequest) (string, error)
	ETHSignTypedData(ctx context.Context, address string, request *ethTypes.SignTypedDataRequest) (string, error)
	ETHSignTransaction(ctx context.Context, address string, request *ethTypes.SignETHTransactionRequest) (string, error)
	ETHSignQuorumPrivateTransaction(ctx context.Context, address string, request *ethTypes.SignQuorumPrivateTransactionRequest) (string, error)
	ETHSignEEATransaction(ctx context.Context, address string, request *ethTypes.SignEEATransactionRequest) (string, error)
	ETHListAccounts(ctx context.Context, namespace string) ([]string, error)
	ETHListNamespaces(ctx context.Context) ([]string, error)
	ETHGetAccount(ctx context.Context, address, namespace string) (*ethTypes.ETHAccountResponse, error)
	ETHVerifySignature(ctx context.Context, request *ethTypes.VerifyPayloadRequest) error
	ETHVerifyTypedDataSignature(ctx context.Context, request *ethTypes.VerifyTypedDataRequest) error
}

type ZKSAccountClient interface {
	ZKSCreateAccount(ctx context.Context, request *zksTypes.CreateZKSAccountRequest) (*zksTypes.ZKSAccountResponse, error)
	ZKSSign(ctx context.Context, address string, request *keymanager.SignPayloadRequest) (string, error)
	ZKSListAccounts(ctx context.Context, namespace string) ([]string, error)
	ZKSListNamespaces(ctx context.Context) ([]string, error)
	ZKSGetAccount(ctx context.Context, address, namespace string) (*zksTypes.ZKSAccountResponse, error)
	ZKSVerifySignature(ctx context.Context, request *zksTypes.VerifyPayloadRequest) error
}

type KeyManagerClient interface {
	Checker() healthz.Check
	EthereumAccountClient
	ZKSAccountClient
}
