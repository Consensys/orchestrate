package registry

import (
	"testing"

	"github.com/containous/traefik/v2/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalProvider(), "Global should have been set")

	testProvider := GlobalProvider()
	Init()
	assert.Equal(t, testProvider, GlobalProvider(), "Provider should not have change after re-initialize")

	var p provider.Provider
	SetGlobalProvider(p)
	assert.Nil(t, GlobalProvider(), "Global should be reset to nil")
}
