package handlers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers/sender"
)

// Init inialize handlers
func Init(ctx context.Context) {
	wg := sync.WaitGroup{}

	// Initialize Jaeger tracer
	wg.Add(1)
	go func() {
		jaeger.Init(ctx)
		wg.Done()
	}()

	// Initialize sender tracer
	wg.Add(1)
	go func() {
		sender.Init(ctx)
		wg.Done()
	}()

	// Wait for all handlers to be ready
	wg.Wait()
}
