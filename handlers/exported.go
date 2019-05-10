package handlers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/store"
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
		store.Init(ctx)
		wg.Done()
	}()

	// Initialize sender tracer
	wg.Add(1)
	go func() {
		producer.Init(ctx)
		wg.Done()
	}()

	// Wait for all handlers to be ready
	wg.Wait()
}
