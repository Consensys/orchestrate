package steps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/envelope"
)

func TestTracker(t *testing.T) {
	tracker := newTracker()

	// Register output on tracker
	ch := make(chan *envelope.Envelope, 10)
	tracker.addOutput("test-output", ch)

	// Input an envelope in channel
	input := &envelope.Envelope{}
	ch <- input

	// Get envelope
	err := tracker.load("test-output", time.Second)
	assert.Nil(t, err, "#1 load should not error")
	assert.Equal(t, input, tracker.current, "#1 envelope should have been loaded")

	// Second load should error
	err = tracker.load("test-output", time.Second)
	assert.NotNil(t, err, "#2 Load should not error")
}
