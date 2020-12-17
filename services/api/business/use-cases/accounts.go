package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

//go:generate mockgen -source=accounts.go -destination=mocks/accounts.go -package=mocks

type AccountUseCases interface {
	GetAccount() GetAccountUseCase
	CreateAccount() CreateAccountUseCase
	UpdateAccount() UpdateAccountUseCase
	SearchAccounts() SearchAccountsUseCase
	FundAccount() FundAccountUseCase
}

type GetAccountUseCase interface {
	Execute(ctx context.Context, address string, tenants []string) (*entities.Account, error)
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, privateKey, chainName, tenantID string) (*entities.Account, error)
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
