package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/examples"
)

// Define a handler method
func handler(txctx *engine.TxContext) {
	txctx.Logger.Infof("Handling %v\n", txctx.In.(examples.Msg))
}

func main() {
	// Instantiate Engine
	cfg := engine.NewConfig()
	eng := engine.NewEngine(&cfg)

	// Register an handler
	eng.Register(handler)

	// Create an input channel of messages
	in := make(chan engine.Msg)

	// Run Engine on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		eng.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- examples.Msg("Message-1")
	in <- examples.Msg("Message-2")
	in <- examples.Msg("Message-3")

	// Close channel & wait for Engine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
	engine.CleanUp()
}
