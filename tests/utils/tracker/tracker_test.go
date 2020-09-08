// +build unit

package tracker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
)

func TestTracker(t *testing.T) {
	tracker := NewTracker()

	// Register output on Tracker
	ch := make(chan *tx.Envelope, 10)
	tracker.AddOutput("test-output", ch)

	// Input an envelope in channel
	input := tx.NewEnvelope()
	ch <- input

	// Get envelope
	err := tracker.Load("test-output", time.Second)
	assert.NoError(t, err, "#1 Load should not error")
	assert.Equal(t, input, tracker.Current, "#1 envelope should have been loaded")

	// Second Load should error
	err = tracker.Load("test-output", time.Second)
	assert.Error(t, err, "#2 Load should not error")
}
