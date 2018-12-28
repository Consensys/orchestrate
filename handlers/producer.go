package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// Producer creates a producer handler
func Producer(p infra.TraceProducer) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		err := p.Produce(ctx.Pb)
		if err != nil {
			ctx.Error(err)
		}
	}
}
