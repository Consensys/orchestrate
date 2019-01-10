package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// Faucet creates a Faucet handler
func Faucet(crediter services.EthCrediter, controller services.EthCreditController) types.HandlerFunc {
	return func(ctx *types.Context) {
		// Interogate credit controller
		amount, ok := controller.ShouldCredit(ctx.T.Chain().ID, *ctx.T.Sender().Address, ctx.T.Tx().Cost())
		if !ok {
			// Credit invalid
			return
		}
		// Credit Valid
		err := crediter.Credit(ctx.T.Chain().ID, *ctx.T.Sender().Address, amount)
		if err != nil {
			// TODO: handle error
			ctx.Error(err)
		}
	}
}
