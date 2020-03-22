package jwt

import (
	"context"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

// Structure to define the parser of the Token and what have to be verify in the Token
type JWT struct {
	conf *Config
}

func New(conf *Config) *JWT {
	return &JWT{conf: conf}
}

const authPrefix = "Bearer "

// Parse and verify the validity of the Token (UUID or Access) and return a struct for a JWT (JSON Web Token)
func (checker *JWT) Check(ctx context.Context) (context.Context, error) {
	if checker.conf.Key == nil {
		// If no KeyFunc provided we deactivate authentication
		return ctx, nil
	}

	// Extract Access Token from context
	bearer, ok := authutils.ParseAuth(authPrefix, authutils.AuthorizationFromContext(ctx))
	if !ok {
		return ctx, errors.UnauthorizedError("missing Access Token")
	}

	// Parse and validate token injected in context
	token, err := checker.conf.Parser.ParseWithClaims(
		bearer,
		&Claims{namespace: checker.conf.ClaimsNamespace},
		checker.conf.Key,
	)
	if err != nil {
		return ctx, errors.UnauthorizedError(err.Error())
	}
	if !token.Valid {
		return ctx, errors.UnauthorizedError("invalid Access Token")
	}

	// Extract multitenancy UUID from token
	tenantID := token.Claims.(*Claims).Orchestrate.TenantID
	if tenantID == "" {
		return ctx, errors.PermissionDeniedError("tenant missing in UUID / Access Token")
	}

	ctx = With(ctx, token)
	ctx = multitenancy.WithTenantID(ctx, tenantID)

	// Enrich context with JWT token
	return ctx, nil
}
