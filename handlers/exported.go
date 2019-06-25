package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers/sender"
)

// Init inialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			jaeger.Init(ctx)
		},

		// Initialize sender tracer
		func() {
			sender.Init(ctx)
		},
	)
}
