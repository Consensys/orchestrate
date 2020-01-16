package grpcserver

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

const AuthorizationHeader = "authorization"

func Auth(auth authentication.Auth, multitenancy bool) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if multitenancy {
			authorization := metautils.ExtractIncoming(ctx).Get(AuthorizationHeader)
			checkedCtx, err := auth.Check(authutils.WithAuthorization(ctx, authorization))
			if err != nil {
				return ctx, err
			}

			// TODO: Uncomment next line after next release of grpc-middleware
			// It is not possible to attach a tag to a go-context in grpc_ctxtags v1.1.0
			// It seems to have been solved on master
			// checkedCtx := grpc_ctxtags.SetInContext(checkedCtx, grpc_ctxtags.Extract(checkedCtx).Set("auth.tenant", tenantID))

			return checkedCtx, nil
		}
		return ctx, nil
	}
}
