package key

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
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
func (checker *Key) Check(ctx context.Context) (*multitenancy.UserInfo, error) {
	if checker == nil || checker.key == "" {
		return nil, nil
	}

	// Extract Key from context
	apiKey := authutils.APIKeyFromContext(ctx)
	if apiKey == "" {
		return nil, nil
	}

	if apiKey != checker.key {
		return nil, errors.UnauthorizedError("invalid API key")
	}

	userInfo := multitenancy.NewAPIKeyUserInfo(apiKey)
	err := userInfo.ImpersonateTenant(authutils.TenantIDFromContext(ctx))
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
