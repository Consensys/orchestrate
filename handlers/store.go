package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// Store creates an handler that can store a trace
func Store(store services.TraceStore) types.HandlerFunc {
	return func(ctx *types.Context) {
		err := store.Store(ctx.T)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
