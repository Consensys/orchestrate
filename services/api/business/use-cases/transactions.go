package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

//go:generate mockgen -source=transactions.go -destination=mocks/transactions.go -package=mocks

/**
Transaction Use Cases
*/
type TransactionUseCases interface {
	SendContractTransaction() SendContractTxUseCase
	SendDeployTransaction() SendDeployTxUseCase
	SendTransaction() SendTxUseCase
	GetTransaction() GetTxUseCase
	SearchTransactions() SearchTransactionsUseCase
}

type GetTxUseCase interface {
	Execute(ctx context.Context, scheduleUUID string, userInfo *multitenancy.UserInfo) (*entities.TxRequest, error)
}

type SearchTransactionsUseCase interface {
	Execute(ctx context.Context, filters *entities.TransactionRequestFilters, userInfo *multitenancy.UserInfo) ([]*entities.TxRequest, error)
}

type SendDeployTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, userInfo *multitenancy.UserInfo) (*entities.TxRequest, error)
}
type SendContractTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, userInfo *multitenancy.UserInfo) (*entities.TxRequest, error)
}

type SendTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, txData string, userInfo *multitenancy.UserInfo) (*entities.TxRequest, error)
}
