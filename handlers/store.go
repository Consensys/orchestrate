package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// CtxStore is used to store context
type CtxStore interface {
	// Send should send raw transaction
	Store(t *types.Trace) error
}

// Store creates an handler that can store a trace
func Store(store CtxStore) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		err := store.Store(ctx.T)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
