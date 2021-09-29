package jwt

import (
	"context"

	"github.com/golang-jwt/jwt"
)

type authCtxKey string

const jwtTokenKey authCtxKey = "token"

// With injects Access Token in context
func With(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, jwtTokenKey, token)
}

// FromContext extracts Access Token from context
func FromContext(ctx context.Context) *jwt.Token {
	token, ok := ctx.Value(jwtTokenKey).(*jwt.Token)
	if !ok {
		return nil
	}
	return token
}
