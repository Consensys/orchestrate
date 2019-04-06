package main

import (
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
	cfg.Slots = 1
	cfg.Partitions = 1
	worker := worker.NewWorker(cfg)

	// Register handlers
	worker.Use(middleware)
	worker.Use(pipeline)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
	go func() { worker.Run(in) }()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"

	// Close channel & wiat for worker to treat all messages
	close(in)
	<-worker.Done()
}
