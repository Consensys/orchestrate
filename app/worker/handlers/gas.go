package handlers

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// GasPricer creates an handler that set Gas Price
func GasPricer(p services.GasPricer) types.HandlerFunc {
	return func(ctx *types.Context) {
		p, err := p.SuggestGasPrice(context.Background(), ctx.T.Chain().ID)
		if err != nil {
			// TODO: handle error
			ctx.Logger.WithError(err).Errorf("gas-pricer: could not suggest gas price")
			ctx.AbortWithError(err)
		} else {
			ctx.T.Tx().SetGasPrice(p)
			// Enrich logger
			ctx.Logger = ctx.Logger.WithFields(log.Fields{
				"tx.gas.price": p.Text(10),
			})
			ctx.Logger.Debugf("gas-pricer: gas price set")
		}
	}
}

// GasEstimator creates an handler that set Gas Limit
func GasEstimator(p services.GasEstimator) types.HandlerFunc {

	pool := &sync.Pool{
		New: func() interface{} { return ethereum.CallMsg{} },
	}

	return func(ctx *types.Context) {
		// Retrieve re-cycled CallMsg
		call := pool.Get().(ethereum.CallMsg)
		defer pool.Put(call)

		// Set CallMsg
		call.From = *ctx.T.Sender().Address
		call.To = ctx.T.Tx().To()
		call.Value = ctx.T.Tx().Value()
		call.Data = ctx.T.Tx().Data()

		g, err := p.EstimateGas(context.Background(), ctx.T.Chain().ID, call)
		if err != nil {
			// TODO: handle error
			ctx.Logger.WithError(err).Errorf("gas-estimator: could not estimate gas limit")
			ctx.AbortWithError(err)
		} else {
			// Set gas limit on context
			ctx.T.Tx().SetGasLimit(g)

			// Enrich logger
			ctx.Logger = ctx.Logger.WithFields(log.Fields{
				"tx.gas.limit": g,
			})
			ctx.Logger.Debugf("gas-estimator: gas limit set")
		}
		
	}
}
