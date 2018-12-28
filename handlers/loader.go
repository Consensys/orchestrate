package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
)

// Loader creates an handler loading input
func Loader(u infra.Unmarshaller) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// Unmarshal message
		err := u.Unmarshal(ctx.Msg, ctx.Pb)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}

		// Load Trace from protobuffer
		protobuf.LoadTrace(ctx.Pb, ctx.T)
	}
}
