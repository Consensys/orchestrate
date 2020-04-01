// +build unit

package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalChainRegistry(), "Global should have been set")

	var chainRegistry *ChainRegistry
	SetGlobalChainRegistry(chainRegistry)
	assert.Nil(t, GlobalChainRegistry(), "Global should be reset to nil")
}
