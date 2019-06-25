package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/handlers/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/handlers/producer"
)

// Init initialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			jaeger.Init(ctx)
		},
		// Initialize Nonce manager
		func() {
			nonce.Init(ctx)
		},
		// Initialize PrepareMsg
		func() {
			producer.Init(ctx)
		},
	)
}
