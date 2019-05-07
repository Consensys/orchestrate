package handlers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/signer"
)

// Init inialize handlers
func Init(ctx context.Context) {
	wg := sync.WaitGroup{}

	// Initialize signer
	wg.Add(1)
	go func() {
		signer.Init(ctx)
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
