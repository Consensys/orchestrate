package jwt

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
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
func (checker *JWT) Check(ctx context.Context) (*multitenancy.UserInfo, error) {
	// Extract Access Token from context
	bearerToken, ok := authutils.ParseBearerToken(authutils.AuthorizationFromContext(ctx))
	if !ok {
		return nil, nil
	}

	// Parse and validate token injected in context
	claims, err := checker.validator.ValidateToken(ctx, bearerToken)
	if err != nil {
		return nil, errors.UnauthorizedError(err.Error())
	}

	userInfo := multitenancy.NewJWTUserInfo(claims, authutils.AuthorizationFromContext(ctx))

	// Impersonate tenant
	err = userInfo.ImpersonateTenant(authutils.TenantIDFromContext(ctx))
	if err != nil {
		return nil, err
	}

	// Impersonate username
	err = userInfo.ImpersonateUsername(authutils.UsernameFromContext(ctx))
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
