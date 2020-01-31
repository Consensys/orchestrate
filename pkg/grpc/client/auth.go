package grpcclient

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

type PerRPCCredentials struct{}

func (cred *PerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	headers := make(map[string]string)
	auth := authutils.AuthorizationFromContext(ctx)
	if auth != "" {
		headers[authentication.AuthorizationHeader] = auth
	}
	apiKey := authutils.APIKeyFromContext(ctx)
	if apiKey != "" {
		headers[authentication.APIKeyHeader] = apiKey
	}
	return headers, nil
}

func (cred *PerRPCCredentials) RequireTransportSecurity() bool {
	return false
}
