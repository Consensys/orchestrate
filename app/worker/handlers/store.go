package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// TraceLoader creates and handler that load traces
func TraceLoader(store infra.TraceStore) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		_, _, err := store.LoadByTxHash(context.Background(), ctx.T.GetChain().GetId(), ctx.T.GetReceipt().GetTxHash(), ctx.T)
		if err != nil {
			// We got an error, possibly due to timeout connexion to database or something else
			// TODO: what should we do in case of error?
			ctx.Error(err)
		}
	}
}
