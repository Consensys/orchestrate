package crafter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalCrafter(), "Global should have been set")

	var c Crafter
	SetGlobalCrafter(c)
	assert.Nil(t, GlobalCrafter(), "Global should be reset to nil")
}
