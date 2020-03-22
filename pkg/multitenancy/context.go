package multitenancy

import (
	"context"
)

type tenancyCtxKey string

const TenantIDKey tenancyCtxKey = "tenant_id"
const DefaultTenantIDName = "_"

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

func TenantIDFromContext(ctx context.Context) string {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	if !ok {
		return DefaultTenantIDName
	}
	return tenantID
}
