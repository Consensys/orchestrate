package utils

import (
	"context"
)

type ethclientKey string

func RetryConnectionError(ctx context.Context, flag bool) context.Context {
	return context.WithValue(ctx, ethclientKey("retry.connection-err"), flag)
}

func ShouldRetryConnectionError(ctx context.Context) bool {
	flag, ok := ctx.Value(ethclientKey("retry.connection-err")).(bool)
	return ok && flag
}
