package main

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Define a pipeline handler
func pipeline(ctx *types.Context) {
	fmt.Printf("* Pipeline handling %v\n", ctx.Msg.(string))
}

// Define a middleware handler
func middleware(ctx *types.Context) {
	// Start middleware execution
	fmt.Printf("* Middleware starts handling %v\n", ctx.Msg.(string))

	// Trigger execution of pending handlers
	ctx.Next()

	// Executed after pending handlers have executed
	fmt.Printf("* Middleware finishes handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate worker (limited to 1 message processed at a time)
	worker := core.NewWorker(1)

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
