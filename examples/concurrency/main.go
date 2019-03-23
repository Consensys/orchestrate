package main

import (
	"context"
	"fmt"
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
	// Instantiate worker that can treat 100 message concurrently in 100 distinct partitions
	cfg := worker.NewConfig()
	cfg.Slots = 100
	cfg.Partitions = 100
	worker := worker.NewWorker(context.Background(), cfg)

	// Register handler
	h := ExampleHandler{0, 0}
	worker.Partitionner(func(msg interface{}) []byte { return []byte(msg.(string)) })
	worker.Use(h.handleSafe)
	worker.Use(h.handleUnsafe)

	// Start worker
	in := make(chan interface{})
	go func() { worker.Run(in) }()

	// Feed 10000 to the worker
	for i := 0; i < 10000; i++ {
		in <- fmt.Sprintf("%v-%v", "Message", i)
	}

	// Close channel
	close(in)
	<-worker.Done()

	// Print counters
	fmt.Printf("* Safe counter: %v\n", h.safeCounter)
	fmt.Printf("* Unsafe counter: %v\n", h.unsafeCounter)
}
