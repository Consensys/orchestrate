// +build unit

package configwatcher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottle(t *testing.T) {
	in := make(chan interface{}, 3)
	out := make(chan interface{}, 3)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		Throttle(ctx, 100*time.Millisecond, in, out)
		close(done)
	}()

	// Send 3 messages to input channel
	in <- "msg1"
	time.Sleep(20 * time.Millisecond)
	in <- "msg2"
	in <- "msg3"

	// Receive message from output channel and ensure 1st message has been dropped
	msg := <-out
	assert.Equal(t, "msg1", msg.(string), "Output message should be correct")

	msg = <-out
	assert.Equal(t, "msg3", msg.(string), "Output message should be correct")
	select {
	case <-out:
		t.Errorf("Output channel should have received only one input")
	default:
	}

	cancel()
	<-done
}
