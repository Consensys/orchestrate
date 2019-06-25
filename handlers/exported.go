package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/faucet"
	gasestimator "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/gas-estimator"
	gaspricer "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/gas-pricer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers/producer"
)

// Init inialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize crafter
		func() {
			crafter.Init(ctx)
		},

		// Initialize faucet
		func() {
			faucet.Init(ctx)
		},

		// Initialize Gas Estimator
		func() {
			gasestimator.Init(ctx)
		},

		// Initialize Gas Pricer
		func() {
			gaspricer.Init(ctx)
		},

		// Initialize Producer
		func() {
			producer.Init(ctx)
		},
	)
}
