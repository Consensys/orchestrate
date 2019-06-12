package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/examples"
)

// Define a pipeline handler
func pipeline(txctx *engine.TxContext) {
	txctx.Logger.Infof("Pipeline handling %v\n", txctx.Msg.(examples.Msg))
}

// Define a middleware handler
func middleware(txctx *engine.TxContext) {
	// Start middleware execution
	txctx.Logger.Infof("Middleware starts handling %v\n", txctx.Msg.(examples.Msg))

	// Trigger execution of pending handlers
	txctx.Next()

	// Executed after pending handlers have executed
	txctx.Logger.Infof("Middleware finishes handling %v\n", txctx.Msg.(examples.Msg))
}

func main() {
	cfg := engine.NewConfig()
	myEngine := engine.NewEngine(&cfg)

	// Register handlers
	myEngine.Register(middleware)
	myEngine.Register(pipeline)

	// Create an input channel of messages
	in := make(chan engine.Msg)

	// Run myEngine on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		myEngine.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- examples.Msg("Message-1")
	in <- examples.Msg("Message-2")

	// Close channel & wait for myEngine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
	engine.CleanUp()
}
