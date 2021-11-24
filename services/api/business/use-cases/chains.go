package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

//go:generate mockgen -source=chains.go -destination=mocks/chains.go -package=mocks

type ChainUseCases interface {
	RegisterChain() RegisterChainUseCase
	GetChain() GetChainUseCase
	SearchChains() SearchChainsUseCase
	UpdateChain() UpdateChainUseCase
	DeleteChain() DeleteChainUseCase
}

type RegisterChainUseCase interface {
	Execute(ctx context.Context, chain *entities.Chain, fromLatest bool, userInfo *multitenancy.UserInfo) (*entities.Chain, error)
}

type GetChainUseCase interface {
	Execute(ctx context.Context, uuid string, userInfo *multitenancy.UserInfo) (*entities.Chain, error)
}

type SearchChainsUseCase interface {
	Execute(ctx context.Context, filters *entities.ChainFilters, userInfo *multitenancy.UserInfo) ([]*entities.Chain, error)
}

type UpdateChainUseCase interface {
	Execute(ctx context.Context, chain *entities.Chain, userInfo *multitenancy.UserInfo) (*entities.Chain, error)
}

type DeleteChainUseCase interface {
	Execute(ctx context.Context, uuid string, userInfo *multitenancy.UserInfo) error
}
