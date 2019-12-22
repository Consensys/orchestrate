package dynamic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	node := &Node{}
	node.SetDefault()
	assert.NotNil(t, node.Listener, "Listener should be set")
	assert.NotEqual(t, "", node.ID, "Listener should be set")
}
