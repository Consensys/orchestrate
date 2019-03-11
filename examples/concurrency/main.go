package main

import (
	"fmt"
	"sync/atomic"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// ExampleHandler is an handler that increment counters
type ExampleHandler struct {
	safeCounter   uint32
	unsafeCounter uint32
}

func (h *ExampleHandler) handleSafe(ctx *types.Context) {
	// Increment counter using atomic
	atomic.AddUint32(&h.safeCounter, 1)
}

func (h *ExampleHandler) handleUnsafe(ctx *types.Context) {
	// Increment counter with no concurrent protection
	h.unsafeCounter++
}

func main() {
	// Instantiate worker (that can treat 1000 message in parallel)
	worker := core.NewWorker(1000)

	// Register handler
	h := ExampleHandler{0, 0}
	worker.Use(h.handleSafe)
	worker.Use(h.handleUnsafe)

	// Start worker
	in := make(chan interface{})
	go func() { worker.Run(in) }()

	// Feed 10000 to the worker
	for i := 0; i < 10000; i++ {
		in <- "Message"
	}

	// Close channel
	close(in)
	<-worker.Done()

	// Print counters
	fmt.Printf("* Safe counter: %v\n", h.safeCounter)
	fmt.Printf("* Unsafe counter: %v\n", h.unsafeCounter)
}
