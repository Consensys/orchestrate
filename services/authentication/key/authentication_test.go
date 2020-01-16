package key

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

func TestAuth(t *testing.T) {
	a := NewAuth("test-key")

	ctx := authutils.WithAuthorization(context.Background(), "APIKey test-key")
	_, err := a.Check(ctx)
	assert.NoError(t, err, "#1 Check should be valid")

	ctx = authutils.WithAuthorization(context.Background(), "APIKey test-key-invalid")
	_, err = a.Check(ctx)
	assert.Error(t, err, "#2 Check should be invalid")

	ctx = authutils.WithAuthorization(context.Background(), "Bearer test-key")
	_, err = a.Check(ctx)
	assert.Error(t, err, "#3 Check should be invalid")
}
