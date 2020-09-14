package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
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
	Execute(ctx context.Context, txRequestUUID string, tenants []string) (*entities.TxRequest, error)
}

type SearchTransactionsUseCase interface {
	Execute(ctx context.Context, filters *entities.TransactionFilters, tenants []string) ([]*entities.TxRequest, error)
}

type SendDeployTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error)
}
type SendContractTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error)
}

type SendTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, txData, tenantID string) (*entities.TxRequest, error)
}
