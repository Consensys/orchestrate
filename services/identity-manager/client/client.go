package client

import (
	"context"

	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type AccountClient interface {
	CreateAccount(ctx context.Context, request *types.CreateAccountRequest) (*types.AccountResponse, error)
	SearchAccounts(ctx context.Context, filters *entities.AccountFilters) ([]*types.AccountResponse, error)
	GetAccount(ctx context.Context, address string) (*types.AccountResponse, error)
	ImportAccount(ctx context.Context, request *types.ImportAccountRequest) (*types.AccountResponse, error)
	UpdateAccount(ctx context.Context, address string, request *types.UpdateAccountRequest) (*types.AccountResponse, error)
	SignPayload(ctx context.Context, address string, request *types.SignPayloadRequest) (string, error)
}

type IdentityManagerClient interface {
	Checker() healthz.Check
	AccountClient
}
