package authentication

import (
	"context"
	"net/textproto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// The TenantID Header have to be used only between tx-listener and envelope-store
const TenantIDHeader = "X-Tenant-ID"
const APIKeyHeader = "X-API-Key"
const AuthorizationHeader = "Authorization"

type AuthFunc func(ctx context.Context) (context.Context, error)

type Auth interface {
	Check(ctx context.Context) (context.Context, error)
}

type combinedAuth struct {
	auths []Auth
}

func (a *combinedAuth) Check(ctx context.Context) (context.Context, error) {
	var err error
	for _, auth := range a.auths {
		ctx, err = auth.Check(ctx)
		if err == nil {
			return ctx, nil
		}
	}
	return ctx, err
}

func CombineAuth(auths ...Auth) Auth {
	return &combinedAuth{auths: auths}
}

// CredMatcher Verifies that the header is part of the accepted headers
func CredMatcher(headerKey string) (mdName string, ok bool) {
	headerKey = textproto.CanonicalMIMEHeaderKey(headerKey)
	if headerKey == textproto.CanonicalMIMEHeaderKey(TenantIDHeader) ||
		headerKey == textproto.CanonicalMIMEHeaderKey(APIKeyHeader) ||
		headerKey == textproto.CanonicalMIMEHeaderKey(AuthorizationHeader) {
		return headerKey, true
	}

	return runtime.DefaultHeaderMatcher(headerKey)
}
