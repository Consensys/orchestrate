package utils

import (
	"context"
)

type authCtxKey string

const authorizationKey authCtxKey = "authorization"
const apiKey authCtxKey = "api-key"
const tenantID authCtxKey = "x-tenant-id"
const username authCtxKey = "x-username"

func WithAuthorization(ctx context.Context, authorization string) context.Context {
	return context.WithValue(ctx, authorizationKey, authorization)
}

func AuthorizationFromContext(ctx context.Context) string {
	authorization, _ := ctx.Value(authorizationKey).(string)
	return authorization
}

func WithAPIKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, apiKey, key)
}

func APIKeyFromContext(ctx context.Context) string {
	authorization, _ := ctx.Value(apiKey).(string)
	return authorization
}

func WithTenantID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, tenantID, value)
}

func TenantIDFromContext(ctx context.Context) string {
	tenant, _ := ctx.Value(tenantID).(string)
	return tenant
}

func WithUsername(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, username, value)
}

func UsernameFromContext(ctx context.Context) string {
	username, _ := ctx.Value(username).(string)
	return username
}
