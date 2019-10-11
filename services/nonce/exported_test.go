package nonce

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, m, "Faucet should have been set")

	var manager Manager
	SetGlobalManager(manager)
	assert.Nil(t, GlobalManager(), "Global should be reset to nil")
}
