// +build unit

package key

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	authutils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	c := New("test-key")

	ctx := authutils.WithAPIKey(context.Background(), "test-key")
	ctx = multitenancy.WithTenantID(ctx, "")
	checkedCtx, err := c.Check(ctx)
	assert.NoError(t, err, "#1 Check should be valid")
	assert.Equal(t, multitenancy.DefaultTenant, multitenancy.TenantIDFromContext(checkedCtx), "#1 Impersonated tenant should be valid")
	assert.Equal(t, []string{multitenancy.Wildcard}, multitenancy.AllowedTenantsFromContext(checkedCtx), "#1 Allowed tenants should be valid")

	ctx = authutils.WithAPIKey(context.Background(), "test-key-invalid")
	_, err = c.Check(ctx)
	assert.Error(t, err, "#2 Check should be invalid")

	ctx = authutils.WithAPIKey(context.Background(), "Bearer test-key")
	_, err = c.Check(ctx)
	assert.Error(t, err, "#3 Check should be invalid")

	ctx = authutils.WithAPIKey(context.Background(), "test-key")
	ctx = multitenancy.WithTenantID(ctx, "foo")
	checkedCtx, err = c.Check(ctx)
	assert.NoError(t, err, "#4 Check should be valid")
	assert.Equal(t, "foo", multitenancy.TenantIDFromContext(checkedCtx), "#4 Impersonated tenant should be valid")
	assert.Equal(t, []string{"foo", multitenancy.DefaultTenant}, multitenancy.AllowedTenantsFromContext(checkedCtx), "#4 Allowed tenants should be valid")

}
