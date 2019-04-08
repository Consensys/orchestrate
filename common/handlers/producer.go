package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Producer creates a producer handler
func Producer(p services.Producer) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Produce Envelope
		err := p.Produce(txctx.Envelope)
		if err != nil {
			txctx.Error(err)
		}
	}
}
