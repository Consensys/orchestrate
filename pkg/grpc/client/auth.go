package grpcclient

import (
	"context"

	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

type PerRPCCredentials struct{}

func (cred *PerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	headers := make(map[string]string)
	auth := authutils.AuthorizationFromContext(ctx)
	if auth != "" {
		headers[grpcserver.AuthorizationHeader] = auth
	}
	return headers, nil
}

func (cred *PerRPCCredentials) RequireTransportSecurity() bool {
	return false
}
