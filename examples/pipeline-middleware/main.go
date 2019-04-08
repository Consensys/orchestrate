package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Define a pipeline handler
func pipeline(ctx *engine.TxContext) {
	ctx.Logger.Infof("Pipeline handling %v\n", ctx.Msg.(string))
}

// Define a middleware handler
func middleware(ctx *engine.TxContext) {
	// Start middleware execution
	ctx.Logger.Infof("Middleware starts handling %v\n", ctx.Msg.(string))

	// Trigger execution of pending handlers
	ctx.Next()

	// Executed after pending handlers have executed
	ctx.Logger.Infof("Middleware finishes handling %v\n", ctx.Msg.(string))
}

func main() {
	cfg := engine.NewConfig()
	engine := engine.NewEngine(&cfg)

	// Register handlers
	engine.Use(middleware)
	engine.Use(pipeline)

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
