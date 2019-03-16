package handlers

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	coreWorker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Faucet creates a Faucet handler
func Faucet(faucet services.Faucet, creditAmount *big.Int) coreWorker.HandlerFunc {
	return func(ctx *coreWorker.Context) {
		faucetRequest := &services.FaucetRequest{
			ChainID: ctx.T.Chain.ID(),
			Address: ctx.T.Sender.Address(),
			Value:   creditAmount,
		}
		amount, approved, err := faucet.Credit(context.Background(), faucetRequest)
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
