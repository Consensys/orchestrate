package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// NonceHandler creates and return an handler for nonce
func NonceHandler(m services.NonceManager) types.HandlerFunc {
	return func(ctx *types.Context) {
		// Retrieve chainID and sender address
		chainID, a := ctx.T.Chain().ID, ctx.T.Sender().Address

		// Retrieve locked nonce from manager
		n, err := m.Obtain(chainID, *a)
		if err != nil {
			ctx.AbortWithError(err)
			return
		}
		err = n.Lock()
		if err != nil {
			ctx.AbortWithError(err)
			return
		}
		defer n.Unlock()

		// Set Nonce value on Trace
		v, err := n.Get()
		if err != nil {
			ctx.AbortWithError(err)
			return
		}
		ctx.T.Tx().SetNonce(v)

		// Execute pending handlers (note that we do not release lock while executing pending handlers)
		ctx.Next()

		// Increment nonce in Manager
		// TODO: we should ensure pending handlers have correctly executed before incrementing
		err = n.Set(v + 1)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
