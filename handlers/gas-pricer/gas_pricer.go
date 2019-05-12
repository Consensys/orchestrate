package gaspricer

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

// Pricer creates an handler that set Gas Price
func Pricer(p ethclient.GasPricer) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetTx().GetTxData().GetGasPrice() == "" {
			// Request a gas price suggestion
			p, err := p.SuggestGasPrice(txctx.Context(), txctx.Envelope.Chain.ID())
			if err != nil {
				// TODO: handle error
				txctx.Logger.WithError(err).Errorf("gas-pricer: could not suggest gas price")
				_ = txctx.AbortWithError(err)
			} else {
				txctx.Envelope.Tx.TxData.SetGasPrice(p)
				// Enrich logger
				txctx.Logger = txctx.Logger.WithFields(log.Fields{
					"tx.gas.price": p.Text(10),
				})
				txctx.Logger.Debugf("gas-pricer: gas price set")
			}
		} else {
			// Gas price has already been set
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"tx.gas.price": txctx.Envelope.GetTx().GetTxData().GetGasPrice(),
			})
		}
	}
}
