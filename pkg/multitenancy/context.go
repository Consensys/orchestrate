package multitenancy

import (
	"context"
)

type tenancyCtxKey string

const (
	TenantIDKey       tenancyCtxKey = "tenant_id"
	AllowedTenantsKey tenancyCtxKey = "allowed_tenants"
	DefaultTenant                   = "_"
)

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

func TenantIDValue(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	return tenantID, ok
}

func TenantIDFromContext(ctx context.Context) string {
	tenantID, ok := TenantIDValue(ctx)
	if ok {
		return tenantID
	}

	return DefaultTenant
}

func WithAllowedTenants(ctx context.Context, tenants []string) context.Context {
	return context.WithValue(ctx, AllowedTenantsKey, tenants)
}

func AllowedTenantsFromContext(ctx context.Context) []string {
	tenants, ok := ctx.Value(AllowedTenantsKey).([]string)
	if ok {
		return tenants
	}
	return []string{DefaultTenant}
}
