package handlers

import (
	"context"
	"sync"

	registryClient "gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/handlers/dispatcher"
)

// Init inialize handlers
func Init(ctx context.Context) {
	wg := sync.WaitGroup{}

	wg.Add(2)
	// Initialize Producer
	go func() {
		dispatcher.Init(ctx)
		wg.Done()
	}()
	// Initialize the registryClient
	go func() {
		registryClient.Init(ctx)
		wg.Done()
	}()

	// Wait for all handlers to be ready
	wg.Wait()
}
