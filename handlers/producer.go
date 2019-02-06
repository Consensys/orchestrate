package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Producer creates a producer handler
func Producer(p services.Producer) types.HandlerFunc {
	return func(ctx *types.Context) {
		// Produce trace protobuffer
		err := p.Produce(ctx.T)
		if err != nil {
			ctx.Error(err)
		}
	}
}
