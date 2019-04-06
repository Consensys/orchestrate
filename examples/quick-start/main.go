package main

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Define a handler method
func handler(ctx *worker.Context) {
	ctx.Logger.Infof("Handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate worker with 1 partition to treat messages
	cfg := worker.NewConfig()
	cfg.Partitions = 1
	worker := worker.NewWorker(cfg)

	// Register handler
	worker.Use(handler)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
	go func() { worker.Run(in) }()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"
	in <- "Message-3"

	// Close channel & wiat for worker to treat all messages
	close(in)
	<-worker.Done()
}
