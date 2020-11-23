package key

import (
	"context"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
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
	if checker.key == "" {
		// If no key provided we deactivate authentication
		return ctx, nil
	}

	// Extract Key from context
	apiKey := authutils.APIKeyFromContext(ctx)
	if apiKey == "" {
		return ctx, errors.UnauthorizedError("missing API key")
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
