package key

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	authutils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/utils"
)

// Key is a Checker for API Key authentication
type Key struct {
	key string
}

func New(key string) *Key {
	return &Key{
		key: key,
	}
}

// Parse and verify the validity of the Token (UUID or Access) and return a struct for a JWT (JSON Web Token)
func (checker *Key) Check(ctx context.Context) (context.Context, error) {
	if checker == nil || checker.key == "" {
		return nil, nil
	}

	// Extract Key from context
	apiKey := authutils.APIKeyFromContext(ctx)
	if apiKey == "" {
		return nil, nil
	}

	if apiKey != checker.key {
		return ctx, errors.UnauthorizedError("invalid API key")
	}

	// Manage multitenancy
	tenantID, err := multitenancy.TenantID(
		multitenancy.Wildcard,
		multitenancy.TenantIDFromContext(ctx),
	)
	if err != nil {
		return ctx, err
	}

	allowedTenants := multitenancy.AllowedTenants(
		multitenancy.Wildcard,
		multitenancy.TenantIDFromContextNoFallback(ctx),
	)

	ctx = multitenancy.WithTenantID(ctx, tenantID)
	ctx = multitenancy.WithAllowedTenants(ctx, allowedTenants)
	return ctx, nil
}
