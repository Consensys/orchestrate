package nodes

import (
	"context"
	"testing"

	"github.com/containous/traefik/v2/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalProvider(), "Global should have been set")

	testProvider := GlobalProvider()
	Init(context.Background())
	assert.Equal(t, testProvider, GlobalProvider(), "Provider should not have change after re-initialize")

	var p provider.Provider
	SetGlobalProvider(p)
	assert.Nil(t, GlobalProvider(), "Global should be reset to nil")
}
