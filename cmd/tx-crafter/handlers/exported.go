package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-crafter/handlers/crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-crafter/handlers/faucet"
	gasestimator "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-crafter/handlers/gas-estimator"
	gaspricer "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-crafter/handlers/gas-pricer"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-crafter/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/handlers/opentracing"
)

type serviceName string

// Init inialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			ctx = context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctx)
		},

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
