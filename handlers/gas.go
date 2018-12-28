package handlers

import (
	"sync"

	"github.com/ethereum/go-ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// GasPricer creates an handler that set Gas Price
func GasPricer(p infra.GasPricer) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		p, err := p.SuggestGasPrice(ctx.T.Chain().ID)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
		}
		ctx.T.Tx().SetGasPrice(p)
	}
}

// GasEstimator creates an handler that set Gas Limit
func GasEstimator(p infra.GasEstimator) infra.HandlerFunc {

	pool := &sync.Pool{
		New: func() interface{} { return ethereum.CallMsg{} },
	}

	return func(ctx *infra.Context) {
		// Retrieve re-cycled CallMsg
		call := pool.Get().(ethereum.CallMsg)
		defer pool.Put(call)

		// Set CallMsg
		call.From = *ctx.T.Sender().Address
		call.To = ctx.T.Tx().To()
		call.Value = ctx.T.Tx().Value()
		call.Data = ctx.T.Tx().Data()

		g, err := p.EstimateGas(ctx.T.Chain().ID, call)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
		}
		ctx.T.Tx().SetGasLimit(g)
	}
}
