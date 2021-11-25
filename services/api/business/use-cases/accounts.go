package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
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
	Execute(ctx context.Context, address ethcommon.Address, userInfo *multitenancy.UserInfo) (*entities.Account, error)
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, privateKey hexutil.Bytes, chainName string, userInfo *multitenancy.UserInfo) (*entities.Account, error)
}

type SearchAccountsUseCase interface {
	Execute(ctx context.Context, filters *entities.AccountFilters, userInfo *multitenancy.UserInfo) ([]*entities.Account, error)
}

type UpdateAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, userInfo *multitenancy.UserInfo) (*entities.Account, error)
}

type FundAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, chainName string, userInfo *multitenancy.UserInfo) error
}
