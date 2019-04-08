package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Define a pipeline handler
func pipeline(txctx *engine.TxContext) {
	txctx.Logger.Infof("Pipeline handling %v\n", txctx.Msg.(string))
}

// Define a middleware handler
func middleware(txctx *engine.TxContext) {
	// Start middleware execution
	txctx.Logger.Infof("Middleware starts handling %v\n", txctx.Msg.(string))

	// Trigger execution of pending handlers
	txctx.Next()

	// Executed after pending handlers have executed
	txctx.Logger.Infof("Middleware finishes handling %v\n", txctx.Msg.(string))
}

func main() {
	cfg := engine.NewConfig()
	engine := engine.NewEngine(&cfg)

	// Register handlers
	engine.Register(middleware)
	engine.Register(pipeline)

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

	// Close channel & wait for Engine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
	engine.CleanUp()
}
