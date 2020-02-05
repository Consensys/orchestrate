package grpcserver

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

func Auth(auth authentication.Auth, multitenancyEnabled bool) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if multitenancyEnabled {
			authorization := metautils.ExtractIncoming(ctx).Get(authentication.AuthorizationHeader)
			apiKey := metautils.ExtractIncoming(ctx).Get(authentication.APIKeyHeader)
			tenantIDFromHeader := metautils.ExtractIncoming(ctx).Get(authentication.TenantIDHeader)

			ctx = authutils.WithAPIKey(authutils.WithAuthorization(multitenancy.WithTenantID(ctx, tenantIDFromHeader), authorization), apiKey)
			checkedCtx, err := auth.Check(ctx)
			if err != nil {
				return ctx, err
			}

			// TODO: Uncomment next line after next release of grpc-middleware
			// It is not possible to attach a tag to a go-context in grpc_ctxtags v1.1.0
			// It seems to have been solved on master
			// checkedCtx := grpc_ctxtags.SetInContext(checkedCtx, grpc_ctxtags.Extract(checkedCtx).Set("auth.tenant", tenantID))

			return checkedCtx, nil
		}
		ctx = multitenancy.WithTenantID(ctx, multitenancy.DefaultTenantIDName)
		return ctx, nil
	}
}
