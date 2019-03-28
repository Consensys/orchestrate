package handlers

import (
	"math/big"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	coreworker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Faucet creates a Faucet handler
func Faucet(faucet services.Faucet, creditAmount *big.Int) coreworker.HandlerFunc {
	return func(ctx *coreworker.Context) {
		faucetRequest := &services.FaucetRequest{
			ChainID: ctx.T.GetChain().ID(),
			Address: ctx.T.GetSender().Address(),
			Value:   creditAmount,
		}
		amount, approved, err := faucet.Credit(ctx.Context(), faucetRequest)
		if err != nil {
			// TODO: handle error
			ctx.Logger.WithError(err).Errorf("faucet: credit error")
			ctx.Error(err)
		} else {
			if !approved {
				ctx.Logger.Debugf("faucet: credit not approved")
			} else {
				ctx.Logger = ctx.Logger.WithFields(log.Fields{
					"faucet.amount": amount.Text(10),
				})
				ctx.Logger.Debugf("faucet: credit approved")
			}
		}
	}
}
