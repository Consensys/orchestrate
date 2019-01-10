package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// Store creates an handler that can store a trace
func Store(store infra.TraceStore) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		err := store.Store(ctx.T)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
