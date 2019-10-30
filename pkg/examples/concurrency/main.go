package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/examples"
)

// ExampleHandler is an handler that increment counters
type ExampleHandler struct {
	safeCounter   uint32
	unsafeCounter uint32
}

func (h *ExampleHandler) handleSafe(txctx *engine.TxContext) {
	// Increment counter using atomic
	atomic.AddUint32(&h.safeCounter, 1)
}

func (h *ExampleHandler) handleUnsafe(txctx *engine.TxContext) {
	// Increment counter with no concurrent protection
	h.unsafeCounter++
}

func main() {
	// Instantiate Engine that can treat 100 message concurrently
	// Instantiate an Engine that can treat 100 message concurrently in 100 distinct partitions
	cfg := engine.NewConfig()
	cfg.Slots = 100
	eng := engine.NewEngine(&cfg)

	// Register handler
	h := ExampleHandler{0, 0}
	eng.Register(h.handleSafe)
	eng.Register(h.handleUnsafe)

	// Run Engine on 100 distinct input channel
	wg := &sync.WaitGroup{}
	inputs := make([]chan engine.Msg, 0)
	for i := 0; i < 100; i++ {
		inputs = append(inputs, make(chan engine.Msg, 100))
		wg.Add(1)
		go func(in chan engine.Msg) {
			eng.Run(context.Background(), in)
			wg.Done()
		}(inputs[i])
	}

	// Feed 10000 to the Engine
	for i := 0; i < 100; i++ {
		for j, in := range inputs {
			in <- examples.Msg(fmt.Sprintf("Message %v-%v", j, i))
		}
	}

	// Close all channels & wait for Engine to treat all messages
	for _, in := range inputs {
		close(in)
	}
	wg.Wait()

	// CleanUp Engine to avoid memory leak
	engine.CleanUp()

	// Print counters
	log.Infof("* Safe counter: %v", h.safeCounter)
	log.Infof("* Unsafe counter: %v", h.unsafeCounter)
}
