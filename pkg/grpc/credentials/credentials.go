package credentials

import (
	"context"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

type PerRPCCredentials struct{}

func (cred *PerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	headers := make(map[string]string)
	authorization := authutils.AuthorizationFromContext(ctx)
	if authorization != "" {
		headers[authutils.AuthorizationHeader] = authorization
	}

	apiKey := authutils.APIKeyFromContext(ctx)
	if apiKey != "" {
		headers[authutils.APIKeyHeader] = apiKey
	}

	tenantID := multitenancy.TenantIDFromContext(ctx)
	if tenantID != "" {
		headers[multitenancy.TenantIDHeader] = tenantID
	}

	return headers, nil
}

func (cred *PerRPCCredentials) RequireTransportSecurity() bool {
	return false
}
