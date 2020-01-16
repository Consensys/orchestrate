package multitenancy

import (
	"context"
)

type tenancyCtxKey string

const TenantIDKey tenancyCtxKey = "tenant_id"

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

func TenantIDFromContext(ctx context.Context) string {
	tenantID, _ := ctx.Value(TenantIDKey).(string)
	return tenantID
}
