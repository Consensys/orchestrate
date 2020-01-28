package jwt

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

// Structure to define the parser of the Token and what have to be verify in the Token
type Auth struct {
	conf *Config
}

func NewAuth(conf *Config) *Auth {
	return &Auth{conf: conf}
}

const authPrefix = "Bearer "

// Parse and verify the validity of the Token (UUID or Access) and return a struct for a JWT (JSON Web Token)
func (a *Auth) Check(ctx context.Context) (context.Context, error) {
	if a.conf.Key == nil {
		// If no KeyFunc provided we deactivate authentication
		return ctx, nil
	}

	// Extract Access Token from context
	auth, ok := authutils.ParseAuth(authPrefix, authutils.AuthorizationFromContext(ctx))
	if !ok {
		return ctx, errors.UnauthorizedError("missing Access Token")
	}

	// Parse and validate token injected in context
	token, err := a.conf.Parser.ParseWithClaims(
		auth,
		&Claims{namespace: a.conf.ClaimsNamespace},
		a.conf.Key,
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
