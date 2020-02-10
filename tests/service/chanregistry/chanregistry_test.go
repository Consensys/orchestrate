package chanregistry

import (
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"github.com/stretchr/testify/assert"
)

func TestChanRegistry(t *testing.T) {
	reg := NewChanRegistry()

	assert.False(t, reg.HasChan("test-key"), "No channel should be registered")
	in := tx.NewBuilder()
	err := reg.Send("test-key", in)
	assert.Error(t, err, "Sending envelope to non registered channel should error")

	// Register channel
	ch := make(chan *tx.Builder, 2)
	reg.Register("test-key", ch)
	assert.True(t, reg.HasChan("test-key"), "Channel should be registered")

	err = reg.Send("test-key", in)
	assert.NoError(t, err, "Sending envelope to registered channel should not error")
	assert.Equal(t, in, <-ch, "Builder should have been sent to channel")
}
