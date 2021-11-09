package jwt

import (
	"context"
)

type contextKey struct{}

// With injects Access Token in context
func With(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, contextKey{}, token)
}

// FromContext extracts Access Token from context
func FromContext(ctx context.Context) string {
	if token, ok := ctx.Value(contextKey{}).(string); ok {
		return token
	}

	return ""
}
