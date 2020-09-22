package utils

import (
	"context"
)

type authCtxKey string

const authorizationKey authCtxKey = "authorization"
const apiKey authCtxKey = "api-Key"

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
