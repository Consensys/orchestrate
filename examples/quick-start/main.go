package main

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Define a handler method
func handler(ctx *engine.TxContext) {
	ctx.Logger.Infof("Handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate Engine with 1 partition to treat messages
	cfg := engine.NewConfig()
	cfg.Partitions = 1
	engine := engine.NewEngine(cfg)

	// Register handler
	engine.Use(handler)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run engine on input channel
	go func() { engine.Run(in) }()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"
	in <- "Message-3"

	// Close channel & wait for engine to treat all messages
	close(in)
	<-engine.Done()
}
