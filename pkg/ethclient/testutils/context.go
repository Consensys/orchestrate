package testutils

import (
	"context"
	"io"
)

type TestCtxKey string

func NewContext(err error, statusCode int, body io.ReadCloser) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, TestCtxKey("resp.error"), err)
	ctx = context.WithValue(ctx, TestCtxKey("resp.statusCode"), statusCode)
	ctx = context.WithValue(ctx, TestCtxKey("resp.body"), body)
	return ctx
}
