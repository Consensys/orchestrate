package handlers

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// GasPricer creates an handler that set Gas Price
func GasPricer(p services.GasPricer) worker.HandlerFunc {
	return func(ctx *worker.Context) {

		if ctx.T.GetTx().GetTxData().GetGasPrice() != "" {
			p, err := p.SuggestGasPrice(ctx.Context(), ctx.T.Chain.ID())
			if err != nil {
				// TODO: handle error
				ctx.Logger.WithError(err).Errorf("gas-pricer: could not suggest gas price")
				ctx.AbortWithError(err)
			} else {
				ctx.T.Tx.TxData.SetGasPrice(p)
				// Enrich logger
				ctx.Logger = ctx.Logger.WithFields(log.Fields{
					"tx.gas.price": p.Text(10),
				})
				ctx.Logger.Debugf("gas-pricer: gas price set")
			}
		}
	}
}

// GasEstimator creates an handler that set Gas Limit
func GasEstimator(p services.GasEstimator) worker.HandlerFunc {

	pool := &sync.Pool{
		New: func() interface{} { return ethereum.CallMsg{} },
	}

	return func(ctx *worker.Context) {

		if ctx.T.GetTx().GetTxData().GetGas() == 0 {
			// Retrieve re-cycled CallMsg
			call := pool.Get().(ethereum.CallMsg)
			defer pool.Put(call)

			To, _ := ctx.T.Tx.GetTxData().ToAddress()
			From, _ := ctx.T.GetSender().Address()
			// Set CallMsg
			call.From = From
			call.To = &To
			call.Value = ctx.T.GetTx().GetTxData().ValueBig()
			call.Data = ctx.T.GetTx().GetTxData().DataBytes()

			g, err := p.EstimateGas(ctx.Context(), ctx.T.GetChain().ID(), call)
			if err != nil {
				// TODO: handle error
				ctx.Logger.WithError(err).Errorf("gas-estimator: could not estimate gas limit")
				ctx.AbortWithError(err)
			} else {
				// Set gas limit on context
				ctx.T.GetTx().GetTxData().SetGas(g)

				// Enrich logger
				ctx.Logger = ctx.Logger.WithFields(log.Fields{
					"tx.gas": g,
				})
				ctx.Logger.Debugf("gas-estimator: gas limit set")
			}
		}
	}
}
