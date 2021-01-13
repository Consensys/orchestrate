// +build unit

package chainregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalManager(), "Global should have been set")

	var mngr *Manager
	SetGlobalManager(mngr)
	assert.Nil(t, GlobalManager(), "Global should be reset to nil")
}
