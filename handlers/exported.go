package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/producer"
)

// Init initialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize decoder
		func() {
			decoder.Init(ctx)
		},
		// Initialize Producer
		func() {
			producer.Init(ctx)
		},
	)
}
