package types

import (
	"math/rand"
	"sync"
	"testing"
	"time"
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
	w := NewWorker(100)
	w.Use(h.Handler(t))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- "test"
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if len(h.handled) != rounds {
		t.Errorf("Worker: expected %v rounds but got %v", rounds, len(h.handled))
	}
}
