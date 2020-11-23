package grpcauth

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

func Auth(checker auth.Checker, multitenancyEnabled bool) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if multitenancyEnabled {
			metadata := metautils.ExtractIncoming(ctx)

			authorization := metadata.Get(authutils.AuthorizationHeader)
			apiKey := metadata.Get(authutils.APIKeyHeader)
			tenantIDFromHeader := metadata.Get(multitenancy.TenantIDHeader)

			ctx = authutils.WithAPIKey(authutils.WithAuthorization(multitenancy.WithTenantID(ctx, tenantIDFromHeader), authorization), apiKey)
			checkedCtx, err := checker.Check(ctx)
			if err != nil {
				return ctx, err
			}

			// TODO: Uncomment next line after next release of grpc-middleware
			// It is not possible to attach a tag to a go-context in grpc_ctxtags v1.1.0
			// It seems to have been solved on master
			// checkedCtx := grpc_ctxtags.SetInContext(checkedCtx, grpc_ctxtags.Extract(checkedCtx).Set("auth.tenant", tenantID))

			return checkedCtx, nil
		}
		ctx = multitenancy.WithTenantID(ctx, multitenancy.DefaultTenant)
		return ctx, nil
	}
}
