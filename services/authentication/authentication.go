package authentication

import "context"

// The TenantID Header have to be used only between tx-listener  and envelop-store
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
