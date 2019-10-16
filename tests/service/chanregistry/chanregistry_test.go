package chanregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/envelope"
)

func TestChanRegistry(t *testing.T) {
	reg := NewChanRegistry()

	assert.False(t, reg.HasChan("test-key"), "No channel should be registered")
	in := &envelope.Envelope{}
	err := reg.Send("test-key", in)
	assert.NotNil(t, err, "Sending envelope to non registered channel should error")

	// Register channel
	ch := make(chan *envelope.Envelope, 2)
	reg.Register("test-key", ch)
	assert.True(t, reg.HasChan("test-key"), "Channel should be registered")

	err = reg.Send("test-key", in)
	assert.Nil(t, err, "Sending envelope to registered channel should not error")
	assert.Equal(t, in, <-ch, "Envelope should have been sent to channel")
}
