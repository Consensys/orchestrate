package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
)

// Producer creates a producer handler
func Producer(p services.Producer) core.HandlerFunc {
	return func(ctx *core.Context) {
		// Produce trace protobuffer
		err := p.Produce(ctx.T)
		if err != nil {
			ctx.Error(err)
		}
	}
}
