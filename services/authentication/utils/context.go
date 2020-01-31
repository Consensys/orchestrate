package utils

import (
	"context"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
)

type authCtxKey string
type apiCtxKey string

const authorizationKey authCtxKey = "authorization"
const apiKey apiCtxKey = "api-Key"

func WithAuthorization(ctx context.Context, authorization string) context.Context {
	return context.WithValue(ctx, authorizationKey, authorization)
}

func AuthorizationFromContext(ctx context.Context) string {
	auth, _ := ctx.Value(authorizationKey).(string)
	return auth
}

func WithAPIKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, apiKey, key)
}

func APIKeyFromContext(ctx context.Context) string {
	auth, _ := ctx.Value(apiKey).(string)
	return auth
}

func AddAuthorizationHeader(req *http.Request) {
	auth := AuthorizationFromContext(req.Context())
	if auth != "" {
		req.Header.Add(authentication.AuthorizationHeader, auth)
	}
}

func AddXAPIKeyHeader(req *http.Request) {
	apiKey := APIKeyFromContext(req.Context())
	if apiKey != "" {
		req.Header.Add(authentication.APIKeyHeader, apiKey)
	}
}
