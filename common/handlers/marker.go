package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Marker creates an handler that mark offsets
func Marker(offset services.OffsetMarker) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Marker is expected to be registered as one of the firt handlers so we are sure we alway mark messages
		txctx.Next()

		// Mark message
		err := offset.Mark(txctx.Msg)
		if err != nil {
			txctx.Error(err)
		}
	}
}
