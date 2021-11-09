package jwt

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
)

const (
	authPrefix = "Bearer "
)

type JWT struct {
	validator Validator
}

func New(validator Validator) *JWT {
	return &JWT{
		validator: validator,
	}
}

// Check verifies the jwt token is valid and injects it in the context
func (checker *JWT) Check(ctx context.Context) (context.Context, error) {
	// Extract Access Token from context
	bearerToken, ok := authutils.ParseAuth(authPrefix, authutils.AuthorizationFromContext(ctx))
	if !ok {
		return nil, nil
	}

	// Parse and validate token injected in context
	userClaims, err := checker.validator.ValidateToken(ctx, bearerToken)
	if err != nil {
		return ctx, errors.UnauthorizedError(err.Error())
	}

	ctx = With(ctx, bearerToken)

	// Manage multitenancy
	tenantID, err := multitenancy.TenantID(userClaims.TenantID, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		return ctx, err
	}

	allowedTenants := multitenancy.AllowedTenants(userClaims.TenantID, multitenancy.TenantIDFromContext(ctx))

	ctx = multitenancy.WithTenantID(ctx, tenantID)
	ctx = multitenancy.WithAllowedTenants(ctx, allowedTenants)

	return ctx, nil
}
