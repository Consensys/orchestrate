// +build unit

package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalNonceManager(), "Faucet should have been set")

	var n *NonceManager
	SetGlobalNonceManager(n)
	assert.Nil(t, GlobalNonceManager(), "Global should be reset to nil")
}
