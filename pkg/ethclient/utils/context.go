package utils

import (
	"context"
)

type ethclientKey string

// RetryNotFoundError attach a flag in context to indicate ethclient
// that it should retry if getting a NotFoundError
func RetryNotFoundError(ctx context.Context, flag bool) context.Context {
	return context.WithValue(ctx, ethclientKey("retry.not-found"), flag)
}

// ShouldRetryNotFoundError retrieve flag attached by RetryNotFoundError
func ShouldRetryNotFoundError(ctx context.Context) bool {
	flag, ok := ctx.Value(ethclientKey("retry.not-found")).(bool)
	return ok && flag
}

func RetryConnectionError(ctx context.Context, flag bool) context.Context {
	return context.WithValue(ctx, ethclientKey("retry.connection-err"), flag)
}

func ShouldRetryConnectionError(ctx context.Context) bool {
	flag, ok := ctx.Value(ethclientKey("retry.connection-err")).(bool)
	return ok && flag
}
