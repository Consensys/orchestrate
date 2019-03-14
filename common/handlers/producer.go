package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Producer creates a producer handler
func Producer(p services.Producer) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		// Produce trace protobuffer
		err := p.Produce(ctx.T)
		if err != nil {
			ctx.Error(err)
		}
	}
}
