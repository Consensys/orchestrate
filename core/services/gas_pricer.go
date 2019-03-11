package services

import (
	"context"
	"math/big"
)

// GasPricer is an interfacet to retrieve GasPrice
type GasPricer interface {
	// SuggestGasPrice suggests gas price
	SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error)
}
