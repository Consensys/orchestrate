package client

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/types"
)

//go:generate mockgen -source=client.go -destination=mocks/client.go -package=mocks

type Eth1Client interface {
	CreateEth1Account(ctx context.Context, storeName string, request *types.CreateEth1AccountRequest) (*types.Eth1AccountResponse, error)
	ImportEth1Account(ctx context.Context, storeName string, request *types.ImportEth1AccountRequest) (*types.Eth1AccountResponse, error)
	UpdateEth1Account(ctx context.Context, storeName, address string, request *types.UpdateEth1AccountRequest) (*types.Eth1AccountResponse, error)
	SignEth1(ctx context.Context, storeName, address string, request *types.SignHexPayloadRequest) (string, error)
	SignEth1Data(ctx context.Context, storeName, account string, request *types.SignHexPayloadRequest) (string, error)
	SignTypedData(ctx context.Context, storeName, address string, request *types.SignTypedDataRequest) (string, error)
	SignTransaction(ctx context.Context, storeName, address string, request *types.SignETHTransactionRequest) (string, error)
	SignQuorumPrivateTransaction(ctx context.Context, storeName, address string, request *types.SignQuorumPrivateTransactionRequest) (string, error)
	SignEEATransaction(ctx context.Context, storeName, address string, request *types.SignEEATransactionRequest) (string, error)
	GetEth1Account(ctx context.Context, storeName, address string) (*types.Eth1AccountResponse, error)
	ListEth1Accounts(ctx context.Context, storeName string) ([]string, error)
	DeleteEth1Account(ctx context.Context, storeName, address string) error
	DestroyEth1Account(ctx context.Context, storeName, address string) error
	RestoreEth1Account(ctx context.Context, storeName, address string) error
	ECRecover(ctx context.Context, storeName string, request *types.ECRecoverRequest) (string, error)
	VerifyEth1Signature(ctx context.Context, storeName string, request *types.VerifyEth1SignatureRequest) error
	VerifyTypedDataSignature(ctx context.Context, storeName string, request *types.VerifyTypedDataRequest) error
}

type KeyManagerClient interface {
	Eth1Client
}
