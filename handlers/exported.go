package handlers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/producer"
)

// Init inialize handlers
func Init(ctx context.Context) {
	wg := sync.WaitGroup{}

	// Initialize decoder
	wg.Add(1)
	go func() {
		decoder.Init(ctx)
		wg.Done()
	}()

	// Initialize Producer
	wg.Add(1)
	go func() {
		producer.Init(ctx)
		wg.Done()
	}()

	// Wait for all handlers to be ready
	wg.Wait()
}
