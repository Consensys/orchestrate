package worker

import (
	"context"
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
	w := NewWorker(
		context.Background(),
		Config{Slots: 100, Partitions: 100, Timeout: 60 * time.Second},
	)
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

func TestWorkerStopped(t *testing.T) {
	h := TestHandler{
		mux:     &sync.Mutex{},
		handled: []*Context{},
	}

	// Create new worker and register test handler
	w := NewWorker(
		context.Background(),
		Config{Slots: 100, Partitions: 100, Timeout: 60 * time.Second},
	)
	w.Use(h.Handler(t))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	go func() {
		for i := 1; i <= rounds; i++ {
			in <- "test"
			time.Sleep(time.Millisecond)
		}
		close(in)
	}()

	// Sleep and close
	time.Sleep(300 * time.Millisecond)
	w.Close()

	// Wait for worker to be done
	<-w.Done()

	if len(h.handled) > 500 {
		t.Errorf("Worker: expected max %v rounds but got %v", 500, len(h.handled))
	}

	msgCount := 0
	for range in {
		// We drain messages
		msgCount++
	}

	if len(h.handled)+msgCount != rounds {
		t.Errorf("Worker: expected all %v messages to have been consumed or drained but got consumed=%v drained=%v", rounds, len(h.handled), msgCount)
	}
}
