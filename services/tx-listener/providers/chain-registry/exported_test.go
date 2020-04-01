// +build unit

package chainregistry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalProvider(), "Global should have been set")

	var p *Provider
	SetGlobalProvider(p)
	assert.Nil(t, GlobalProvider(), "Global should be reset to nil")
}
