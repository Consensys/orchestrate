package key

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

// Structure to define the parser of the Token and what have to be verify in the Token
type Auth struct {
	key string
}

func NewAuth(key string) *Auth {
	return &Auth{
		key: key,
	}
}

const apiKeyPrefix = "APIKey "

// Parse and verify the validity of the Token (UUID or Access) and return a struct for a JWT (JSON Web Token)
func (a *Auth) Check(ctx context.Context) (context.Context, error) {
	if a.key == "" {
		// If no key provided we deactivate authentication
		return ctx, nil
	}

	// Extract Key from context
	auth, ok := authutils.ParseAuth(apiKeyPrefix, authutils.AuthorizationFromContext(ctx))
	if !ok {
		return ctx, errors.UnauthorizedError("missing authentication key")
	}

	if auth != a.key {
		return ctx, errors.UnauthorizedError("invalid Token authentication key")
	}

	// Enrich context with JWT token
	return ctx, nil
}
