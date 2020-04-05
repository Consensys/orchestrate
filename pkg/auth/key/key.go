package key

import (
	"context"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
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

	// Grant all privileges to context
	return authutils.GrantAllPrivileges(ctx), nil
}