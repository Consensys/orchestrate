package main

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Define a handler method
func handler(ctx *types.Context) {
	fmt.Printf("* Handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate worker (limited to 1 message processed at a time)
	worker := core.NewWorker(1)

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
