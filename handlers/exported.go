package handlers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/faucet"
	gasestimator "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/gas-estimator"
	gaspricer "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/gas-pricer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/producer"
)

// Init inialize handlers
func Init(ctx context.Context) {
	wg := sync.WaitGroup{}

	// Initialize crafter
	wg.Add(1)
	go func() {
		crafter.Init(ctx)
		wg.Done()
	}()

	// Initialize faucet
	wg.Add(1)
	go func() {
		faucet.Init(ctx)
		wg.Done()
	}()

	// Initialize Gas Estimator
	wg.Add(1)
	go func() {
		gasestimator.Init(ctx)
		wg.Done()
	}()

	// Initialize Gas Pricer
	wg.Add(1)
	go func() {
		gaspricer.Init(ctx)
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
