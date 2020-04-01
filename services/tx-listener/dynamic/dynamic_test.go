// +build unit

package dynamic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	chain := &Chain{}
	chain.SetDefault()
	assert.NotNil(t, chain.Listener, "Listener should be set")
	assert.NotEqual(t, "", chain.UUID, "Listener should be set")
}
