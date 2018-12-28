package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// Marker creates an handler that mark offsets
func Marker(offset infra.OffsetMarker) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		err := offset.Mark(ctx.Msg)
		if err != nil {
			ctx.Error(err)
		}
	}
}
