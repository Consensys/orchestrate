package main

import (
	"fmt"
	"sync/atomic"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// ExampleHandler is an handler that increment counters
type ExampleHandler struct {
	safeCounter   uint32
	unsafeCounter uint32
}

func (h *ExampleHandler) handleSafe(ctx *engine.TxContext) {
	// Increment counter using atomic
	atomic.AddUint32(&h.safeCounter, 1)
}

func (h *ExampleHandler) handleUnsafe(ctx *engine.TxContext) {
	// Increment counter with no concurrent protection
	h.unsafeCounter++
}

func main() {
	// Instantiate an Engine that can treat 100 message concurrently in 100 distinct partitions
	cfg := engine.NewConfig()
	cfg.Slots = 100
	cfg.Partitions = 100
	engine := engine.NewEngine(cfg)

	// Register handler
	h := ExampleHandler{0, 0}
	engine.Partitionner(func(msg interface{}) []byte { return []byte(msg.(string)) })
	engine.Use(h.handleSafe)
	engine.Use(h.handleUnsafe)

	// Start Engine
	in := make(chan interface{})
	go func() { engine.Run(in) }()

	// Feed 10000 to the Engine
	for i := 0; i < 10000; i++ {
		in <- fmt.Sprintf("%v-%v", "Message", i)
	}

	// Close channel
	close(in)
	<-engine.Done()

	// Print counters
	fmt.Printf("* Safe counter: %v\n", h.safeCounter)
	fmt.Printf("* Unsafe counter: %v\n", h.unsafeCounter)
}
