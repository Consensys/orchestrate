package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
)

// Marker creates an handler that mark offsets
func Marker(offset services.OffsetMarker) core.HandlerFunc {
	return func(ctx *core.Context) {
		// Marker is expected to be registered as one of the firt handlers so we are sure we alway mark messages
		ctx.Next()

		// Mark message
		err := offset.Mark(ctx.Msg)
		if err != nil {
			ctx.Error(err)
		}
	}
}
