package handlers

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

// GasPriceEstimator is an interfacet to retrieve GasPrice
type GasPriceEstimator interface {
	// SuggestGasPrice suggests gas price
	SuggestGasPrice(chainID *big.Int) (*big.Int, error)
}

// GasPricer creates an handler that set Gas Price
func GasPricer(p GasPriceEstimator) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		p, err := p.SuggestGasPrice(ctx.T.Chain().ID)
		if err != nil {
			// TODO: handle error
		}
		ctx.T.Tx().SetGasPrice(p)
	}
}

// GasEstimator is an interface to retrieve Gas Cost of a transaction
type GasEstimator interface {
	EstimateGas(chainID *big.Int, call ethereum.CallMsg) (uint64, error)
}

// SimpleGasManager implements methods to get information about Gas using an Ethereum client
type SimpleGasManager struct {
	ec *infra.EthClient
}

// SuggestGasPrice suggests a gas price
func (m *SimpleGasManager) SuggestGasPrice(chainID *big.Int) (*big.Int, error) {
	return m.ec.SuggestGasPrice(context.Background())
}

// EstimateGas suggests a gas limit
func (m *SimpleGasManager) EstimateGas(chainID *big.Int, call ethereum.CallMsg) (uint64, error) {
	return m.ec.EstimateGas(context.Background(), call)
}

// GasLimiter creates an handler that set Gas Limit
func GasLimiter(p GasEstimator) infra.HandlerFunc {

	pool := &sync.Pool{
		New: func() interface{} { return ethereum.CallMsg{} },
	}

	return func(ctx *infra.Context) {
		call := pool.Get().(ethereum.CallMsg)
		defer pool.Put(call)

		call.From = *ctx.T.Sender().Address
		call.To = ctx.T.Tx().To()
		call.Value = ctx.T.Tx().Value()
		call.Data = ctx.T.Tx().Data()

		g, err := p.EstimateGas(ctx.T.Chain().ID, call)
		if err != nil {
			// TODO: handle error
		}
		ctx.T.Tx().SetGasLimit(g)
	}
}
