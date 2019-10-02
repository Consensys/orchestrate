package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/examples"
)

func aborter(txctx *engine.TxContext) {
	txctx.Logger.Infof("Aborting %v\n", txctx.In.(examples.Msg))
	txctx.Abort()
}

// Define a pipeline handler
func pipeline(name string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.Infof("Pipeline-%v handling %v\n", name, txctx.In.(examples.Msg))
	}
}

// Define a middleware handler
func middleware(name string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Start middleware execution
		txctx.Logger.Infof("Middleware-%v starts handling %v\n", name, txctx.In.(examples.Msg))

		// Trigger execution of pending handlers
		txctx.Next()

		// Executed after pending handlers have executed
		txctx.Logger.Infof("Middleware-%v finishes handling %v\n", name, txctx.In.(examples.Msg))
	}
}

func main() {
	// Register handlers
	engine.Init(context.Background())

	// Declare individual handlers
	pipeline1 := pipeline("1")
	pipeline2 := pipeline("2")
	pipeline3 := pipeline("3")
	pipeline4 := pipeline("4")
	middleware1 := middleware("1")
	middleware2 := middleware("2")
	middleware3 := middleware("3")

	// Declare 2 composite handlers
	left := engine.CombineHandlers(middleware2, pipeline2)
	right := engine.CombineHandlers(middleware3, aborter, pipeline3)

	// Declare a forked handler
	fork := func(txctx *engine.TxContext) {
		switch txctx.In.Entrypoint() {
		case "left":
			left(txctx)
		case "right":
			right(txctx)
		}
	}

	// Declare overall composite handler
	handler := engine.CombineHandlers(pipeline1, middleware1, fork, pipeline4)

	// Register composite handler
	engine.Register(handler)

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
	in <- examples.Msg("left")
	in <- examples.Msg("right")

	// Close channel & wait for myEngine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
	engine.CleanUp()
}
