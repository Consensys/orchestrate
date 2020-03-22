package utils

import (
	"context"
)

type authCtxKey string

const authorizationKey authCtxKey = "authorization"
const apiKey authCtxKey = "api-Key"
const allPrivilegesKey authCtxKey = "api-Key"

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

func GrantAllPrivileges(ctx context.Context) context.Context {
	return context.WithValue(ctx, allPrivilegesKey, true)
}

func HasAllPrivileges(ctx context.Context) bool {
	v, _ := ctx.Value(allPrivilegesKey).(bool)
	return v
}
