package client

import (
	"context"

	healthz "github.com/heptiolabs/healthcheck"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type IdentityManagerClient interface {
	Checker() healthz.Check
	CreateIdentity(ctx context.Context, request *types.CreateIdentityRequest) (*types.IdentityResponse, error)
}
