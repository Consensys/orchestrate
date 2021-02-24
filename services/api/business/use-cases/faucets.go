package usecases

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
)

//go:generate mockgen -source=faucets.go -destination=mocks/faucets.go -package=mocks

type FaucetUseCases interface {
	RegisterFaucet() RegisterFaucetUseCase
	UpdateFaucet() UpdateFaucetUseCase
	GetFaucet() GetFaucetUseCase
	SearchFaucets() SearchFaucetsUseCase
	DeleteFaucet() DeleteFaucetUseCase
}

type RegisterFaucetUseCase interface {
	Execute(ctx context.Context, faucet *entities.Faucet) (*entities.Faucet, error)
}

type UpdateFaucetUseCase interface {
	Execute(ctx context.Context, faucet *entities.Faucet, tenants []string) (*entities.Faucet, error)
}

type GetFaucetUseCase interface {
	Execute(ctx context.Context, uuid string, tenants []string) (*entities.Faucet, error)
}

type SearchFaucetsUseCase interface {
	Execute(ctx context.Context, filters *entities.FaucetFilters, tenants []string) ([]*entities.Faucet, error)
}

type DeleteFaucetUseCase interface {
	Execute(ctx context.Context, uuid string, tenants []string) error
}

type GetFaucetCandidateUseCase interface {
	Execute(ctx context.Context, account string, chain *entities.Chain, tenants []string) (*entities.Faucet, error)
}
