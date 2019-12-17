package token

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, auth, "auth Manager should have been set")

	var m authentication.Manager
	SetGlobalAuth(m)
	assert.Nil(t, GlobalAuth(), "Global should be reset to nil")
}
