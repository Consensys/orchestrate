package services

import (
	"math/big"
)

// GasPricer is an interfacet to retrieve GasPrice
type GasPricer interface {
	// SuggestGasPrice suggests gas price
	SuggestGasPrice(chainID *big.Int) (*big.Int, error)
}
