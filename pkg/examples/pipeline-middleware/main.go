package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/examples"
)

func aborter(txctx *engine.TxContext) {
	txctx.Logger.Infof("Aborting %v", txctx.In.(examples.Msg))
	txctx.Abort()
}

// Define a pipeline handler
func pipeline(txctx *engine.TxContext) {
	txctx.Logger.Infof("Pipeline handling %v\n", txctx.In.(examples.Msg))
}

// Define a middleware handler
func middleware(txctx *engine.TxContext) {
	// Start middleware execution
	txctx.Logger.Infof("Middleware starts handling %v\n", txctx.In.(examples.Msg))

	// Trigger execution of pending handlers
	txctx.Next()

	// Executed after pending handlers have executed
	txctx.Logger.Infof("Middleware finishes handling %v\n", txctx.In.(examples.Msg))
}

func main() {
	// Register handlers
	engine.Init(context.Background())
	engine.Register(middleware)
	engine.Register(pipeline)
	engine.Register(aborter)

	// Create an input channel of messages
	in := make(chan engine.Msg)

	// Run myEngine on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		engine.Run(context.Background(), in)
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
