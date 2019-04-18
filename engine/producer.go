package engine

import "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"

// Producer creates a producer handler
func Producer(p services.Producer) HandlerFunc {
	return func(txctx *TxContext) {
		// Produce Envelope
		err := p.Produce(txctx.Envelope)
		if err != nil {
			txctx.Error(err)
		}
	}
}
