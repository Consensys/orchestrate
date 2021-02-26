package httputil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := WithMiddleware(WithService(WithRouter(WithEntryPoint(context.Background(), "test-ep"), "test-router"), "test-service"), "test-middleware")
	assert.Equal(t, "test-ep", EntryPointFromContext(ctx))
	assert.Equal(t, "test-router", RouterFromContext(ctx))
	assert.Equal(t, "test-service", ServiceFromContext(ctx))
	assert.Equal(t, "test-middleware", MiddlewareFromContext(ctx))
}
