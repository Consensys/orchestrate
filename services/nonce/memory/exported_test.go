package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalNonceManager(), "Nonce should have been set")

	var n *NonceManager
	SetGlobalNonceManager(n)
	assert.Nil(t, GlobalNonceManager(), "Global should be reset to nil")
}
