package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Define a handler method
func handler(ctx *engine.TxContext) {
	ctx.Logger.Infof("Handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate worker
	cfg := engine.NewConfig()
	engine := engine.NewEngine(&cfg)

	// Register an handler
	engine.Use(handler)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
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

	// Close channel & wait for worker to treat all messages
	close(in)
	wg.Wait()

	// CleanUp worker to avoid memory leak
	engine.CleanUp()
}
