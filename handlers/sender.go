package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// Sender creates a Sender handler
func Sender(sender infra.TxSender) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		if len(ctx.T.Tx().Raw()) == 0 {
			// Tx is not ready
			// TODO: handle case
			ctx.Abort()
			return
		}

		err := sender.Send(ctx.T.Chain().ID, hexutil.Encode(ctx.T.Tx().Raw()))
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}
