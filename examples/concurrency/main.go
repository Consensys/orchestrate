package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// ExampleHandler is an handler that increment counters
type ExampleHandler struct {
	safeCounter   uint32
	unsafeCounter uint32
}

func (h *ExampleHandler) handleSafe(ctx *worker.Context) {
	// Increment counter using atomic
	atomic.AddUint32(&h.safeCounter, 1)
}

func (h *ExampleHandler) handleUnsafe(ctx *worker.Context) {
	// Increment counter with no concurrent protection
	h.unsafeCounter++
}

func main() {
	// Instantiate worker that can treat 100 message concurrently
	cfg := worker.NewConfig()
	cfg.Slots = 100
	w := worker.NewWorker(&cfg)

	// Register handler
	h := ExampleHandler{0, 0}
	w.Use(h.handleSafe)
	w.Use(h.handleUnsafe)

	// Run worker on 100 distinct input channel
	wg := &sync.WaitGroup{}
	inputs := make([]chan interface{}, 0)
	for i := 0; i < 100; i++ {
		inputs = append(inputs, make(chan interface{}, 100))
		wg.Add(1)
		go func(in chan interface{}) {
			w.Run(context.Background(), in)
			wg.Done()
		}(inputs[i])
	}

	// Feed 10000 to the worker
	for i := 0; i < 100; i++ {
		for j, in := range inputs {
			in <- fmt.Sprintf("Message %v-%v", j, i)
		}
	}

	// Close all channels & wait for worker to treat all messages
	for _, in := range inputs {
		close(in)
	}
	wg.Wait()

	// CleanUp worker to avoid memory leak
	w.CleanUp()

	// Print counters
	fmt.Printf("* Safe counter: %v\n", h.safeCounter)
	fmt.Printf("* Unsafe counter: %v\n", h.unsafeCounter)
}
