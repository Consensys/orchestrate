// +build unit

package key

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
)

func TestKey(t *testing.T) {
	c := New("test-key")

	ctx := authutils.WithAPIKey(context.Background(), "test-key")
	_, err := c.Check(ctx)
	assert.NoError(t, err, "#1 Check should be valid")

	ctx = authutils.WithAPIKey(context.Background(), "test-key-invalid")
	_, err = c.Check(ctx)
	assert.Error(t, err, "#2 Check should be invalid")

	ctx = authutils.WithAPIKey(context.Background(), "Bearer test-key")
	_, err = c.Check(ctx)
	assert.Error(t, err, "#3 Check should be invalid")
}
