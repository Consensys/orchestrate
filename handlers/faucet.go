package handlers

import (
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Faucet creates a Faucet handler
func Faucet(faucet services.Faucet, amountToTransfer *big.Int) types.HandlerFunc {
	return func(ctx *types.Context) {
		faucetRequest := &services.FaucetRequest{
			ChainID: ctx.T.Chain().ID,
			Address: *ctx.T.Sender().Address,
			Value:   amountToTransfer,
		}
		_, _, err := faucet.Credit(faucetRequest)
		if err != nil {
			// TODO: handle error
			ctx.Error(err)
		}
	}
}
