package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

// GasEstimator is an interface to retrieve Gas Cost of a transaction
type GasEstimator interface {
	EstimateGas(ctx context.Context, chainID *big.Int, call ethereum.CallMsg) (uint64, error)
}
