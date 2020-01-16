package utils

import "context"

type authCtxKey string

const authorizationKey authCtxKey = "authorization"

func WithAuthorization(ctx context.Context, authorization string) context.Context {
	return context.WithValue(ctx, authorizationKey, authorization)
}

func AuthorizationFromContext(ctx context.Context) string {
	auth, _ := ctx.Value(authorizationKey).(string)
	return auth
}
