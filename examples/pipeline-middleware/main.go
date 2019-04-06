package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Define a pipeline handler
func pipeline(ctx *worker.Context) {
	ctx.Logger.Infof("Pipeline handling %v\n", ctx.Msg.(string))
}

// Define a middleware handler
func middleware(ctx *worker.Context) {
	// Start middleware execution
	ctx.Logger.Infof("Middleware starts handling %v\n", ctx.Msg.(string))

	// Trigger execution of pending handlers
	ctx.Next()

	// Executed after pending handlers have executed
	ctx.Logger.Infof("Middleware finishes handling %v\n", ctx.Msg.(string))
}

func main() {
	cfg := worker.NewConfig()
	w := worker.NewWorker(cfg)

	// Register handlers
	w.Use(middleware)
	w.Use(pipeline)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		w.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"

	// Close channel & wait for worker to treat all messages
	close(in)
	wg.Wait()

	// CleanUp worker to avoid memory leak
	w.CleanUp()
}
