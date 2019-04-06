package worker

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestHandler struct {
	mux     *sync.Mutex
	handled []*Context
}

func (h *TestHandler) Handler(t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		// We add some randomness in time execution
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}

func TestWorker(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*Context{},
	}

	// Create new worker and register test handler
	w := NewWorker(Config{Slots: 100})
	w.Use(h.Handler(t))

	// Create input channels and prefills it
	ins := make([]chan interface{}, 0)
	for i := 0; i < 50; i++ {
		in := make(chan interface{}, 20)
		for j := 0; j < 20; j++ {
			in <- fmt.Sprintf("test-%v-%v", i, j)
		}
		close(in)
		ins = append(ins, in)
	}

	// Start consuming every input channel
	wg := &sync.WaitGroup{}
	for i := range ins {
		wg.Add(1)
		go func(in <-chan interface{}) {
			w.Run(context.Background(), in)
			wg.Done()
		}(ins[i])
	}

	// Wait for worker to finish consuming
	wg.Wait()

	assert.Len(t, h.handled, 1000, "All messages should have been processed")
}

func TestWorkerStopped(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*Context{},
	}

	// Create new worker and register test handler
	w := NewWorker(
		Config{Slots: 100},
	)
	w.Use(h.Handler(t))

	// Create input channels and prefills it
	ins := make([]chan interface{}, 0)
	for i := 0; i < 50; i++ {
		in := make(chan interface{}, 20)
		for j := 0; j < 20; j++ {
			in <- fmt.Sprintf("test-%v-%v", i, j)
		}
		close(in)
		ins = append(ins, in)
	}

	// Start consuming every input channel
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	for i := range ins {
		wg.Add(1)
		go func(in <-chan interface{}) {
			w.Run(ctx, in)
			wg.Done()
		}(ins[i])
	}

	// Sleep for a short time and interupt
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Wait for worker to finish
	wg.Wait()

	assert.True(t, len(h.handled) < 500, "Expected at least half of the message not to have been consumed")

	// We drain and count all messages that have not been consumed
	count := 0
	for i := range ins {
		for range ins[i] {
			count++
		}
	}

	assert.Equal(t, 1000, len(h.handled)+count, "Expected all message to have either been consumed or still be in input channel")
}

func testSleepingHandler(ctx *Context) {
	time.Sleep(ctx.Keys["duration"].(time.Duration))
}

func makeTimeoutContext(i int) *Context {
	ctx := NewContext()
	ctx.Reset()
	ctx.Prepare([]HandlerFunc{}, log.NewEntry(log.StandardLogger()), nil)

	switch i % 2 {
	case 0:
		ctx.Keys["duration"] = 50 * time.Millisecond
		ctx.Keys["errors"] = 0
	case 1:
		ctx.Keys["duration"] = 100 * time.Millisecond
		ctx.Keys["errors"] = 1
	}
	return ctx
}

func TestTimeoutHandler(t *testing.T) {
	timeoutHandler := TimeoutHandler(testSleepingHandler, 60*time.Millisecond, "Test timeout")

	rounds := 100
	outs := make(chan *Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeTimeoutContext(i)
		go func(ctx *Context) {
			defer wg.Done()
			timeoutHandler(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	assert.Len(t, outs, rounds, "Timeout: processed contexts count should be correct")

	for out := range outs {
		errCount := out.Keys["errors"].(int)
		assert.Len(t, out.T.Errors, errCount, "Timeout: expected correct count of errors")
	}
}
