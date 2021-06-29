package usecases

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:generate mockgen -source=accounts.go -destination=mocks/accounts.go -package=mocks

type AccountUseCases interface {
	GetAccount() GetAccountUseCase
	CreateAccount() CreateAccountUseCase
	UpdateAccount() UpdateAccountUseCase
	SearchAccounts() SearchAccountsUseCase
}

type GetAccountUseCase interface {
	Execute(ctx context.Context, address string, tenants []string) (*entities.Account, error)
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, privateKey hexutil.Bytes, chainName, tenantID string) (*entities.Account, error)
}

type SearchAccountsUseCase interface {
	Execute(ctx context.Context, filters *entities.AccountFilters, tenants []string) ([]*entities.Account, error)
}

type UpdateAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, tenants []string) (*entities.Account, error)
}

type FundAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, chainName string, tenantID string) error
}
