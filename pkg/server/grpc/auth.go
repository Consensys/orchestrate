package grpcserver

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	token_manager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/token"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

// AuthTokenTenant functions used by gRPC interceptor to authenticate the caller with ID / Access Token and extract tenantID
func AuthTokenTenant(ctx context.Context) (context.Context, error) {
	if !viper.GetBool(multitenancy.EnabledViperKey) {
		// Run the next Interceptor
		return ctx, nil
	}
	rawToken, err := grpc_auth.AuthFromMD(ctx, token_manager.HeaderKey)
	if err != nil {
		e := errors.UnauthorizedError("Token Not Found with bearer")
		return nil, e
	}

	token, err := token_manager.GlobalAuth().Verify(rawToken)
	if err != nil {
		e := errors.UnauthorizedError(err.Error())
		return nil, e
	}

	tenantPath := viper.GetString(authentication.TenantNamespaceViperKey)

	tenantIDValue, ok := token.Claims.(jwt.MapClaims)[tenantPath+authentication.TenantIDKey].(string)
	if !ok {
		e := errors.NotFoundError("not able to retrieve the tenant ID: The tenant_id is not present in the ID / Access Token")
		return nil, e
	}

	// Add the Token information and the Tenant Id in the go Context and Tag the Tenant for grpc
	grpc_ctxtags.Extract(ctx).Set("auth.tenant", tenantIDValue)
	newCtx := context.WithValue(context.WithValue(ctx, // nolint:golint // reason
		authentication.TokenInfoKey, token),
		authentication.TenantIDKey, tenantIDValue)

	return newCtx, nil
}
