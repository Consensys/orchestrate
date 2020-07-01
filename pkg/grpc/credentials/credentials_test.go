// +build unit

package credentials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

func TestPerRPCCredentials(t *testing.T) {
	cred := &PerRPCCredentials{}
	h, _ := cred.GetRequestMetadata(context.Background(), "")
	assert.Len(t, h, 1, "Header length should be valid")

	ctx := authutils.WithAuthorization(context.Background(), "test-auth")
	h, _ = cred.GetRequestMetadata(ctx, "")
	require.Len(t, h, 2, "Header length should be valid")
	assert.Equal(t, "test-auth", h[authutils.AuthorizationHeader], "Header should be correct")

	ctx = authutils.WithAPIKey(context.Background(), "test-auth")
	h, _ = cred.GetRequestMetadata(ctx, "")
	require.Len(t, h, 2, "Header length should be valid")
	assert.Equal(t, "test-auth", h[authutils.APIKeyHeader], "Header should be correct")

	ctx = multitenancy.WithTenantID(context.Background(), "test")
	h, _ = cred.GetRequestMetadata(ctx, "")
	require.Len(t, h, 1, "Header length should be valid")
	assert.Equal(t, "test", h[multitenancy.TenantIDHeader], "Header should be correct")
}
