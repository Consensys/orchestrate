package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Define a handler method
func handler(txctx *engine.TxContext) {
	txctx.Logger.Infof("Handling %v\n", txctx.Msg.(string))
}

func main() {
	// Instantiate Engine
	cfg := engine.NewConfig()
	engine := engine.NewEngine(&cfg)

	// Register an handler
	engine.Register(handler)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run Engine on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		engine.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"
	in <- "Message-3"

	// Close channel & wait for Engine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
	engine.CleanUp()
}
