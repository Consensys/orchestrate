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

func TenantIDFromContext(ctx context.Context) string {
	tenantID, _ := ctx.Value(TenantIDKey).(string)
	return tenantID
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
