// +build unit

package multitenancy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTenantID(t *testing.T) {
	assert.Equal(t, "_", TenantIDFromContext(context.Background()))
	ctx := WithTenantID(context.Background(), "test-tenant")
	assert.Equal(t, "test-tenant", TenantIDFromContext(ctx))
}
