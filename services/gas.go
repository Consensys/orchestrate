package services

import (
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
