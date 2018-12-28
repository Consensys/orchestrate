package infra

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

// GasPricer is an interfacet to retrieve GasPrice
type GasPricer interface {
	// SuggestGasPrice suggests gas price
	SuggestGasPrice(chainID *big.Int) (*big.Int, error)
}

// GasEstimator is an interface to retrieve Gas Cost of a transaction
type GasEstimator interface {
	EstimateGas(chainID *big.Int, call ethereum.CallMsg) (uint64, error)
}

// SimpleGasManager implements methods to get information about Gas using an Ethereum client
type SimpleGasManager struct {
	ec *EthClient
}

// NewSimpleGasManager creates a new SimpleGasManager
func NewSimpleGasManager(ec *EthClient) *SimpleGasManager {
	return &SimpleGasManager{ec}
}

// SuggestGasPrice suggests a gas price
func (m *SimpleGasManager) SuggestGasPrice(chainID *big.Int) (*big.Int, error) {
	return m.ec.SuggestGasPrice(context.Background())
}

// EstimateGas suggests a gas limit
func (m *SimpleGasManager) EstimateGas(chainID *big.Int, call ethereum.CallMsg) (uint64, error) {
	return m.ec.EstimateGas(context.Background(), call)
}
