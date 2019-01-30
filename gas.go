package ethereum

import (
	"context"
	"math/big"

	geth "github.com/ethereum/go-ethereum"
)

// GasManagerEthClient is a minimal Ethereum client interface required by a Gas Manager
type GasManagerEthClient interface {
	// Should provide a gas price for a given chain
	SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error)

	// Should provide a gas cost for a transaction on a given chain
	EstimateGas(ctx context.Context, chainID *big.Int, call geth.CallMsg) (uint64, error)
}

// GasManager implements methods to get information about Gas by connecting to an Ethereum client
type GasManager struct {
	ec GasManagerEthClient
}

// NewGasManager creates a new GasManager
func NewGasManager(ec GasManagerEthClient) *GasManager {
	return &GasManager{ec}
}

// SuggestGasPrice suggests a gas price
func (m *GasManager) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	return m.ec.SuggestGasPrice(ctx, chainID)
}

// EstimateGas suggests a gas limit
func (m *GasManager) EstimateGas(ctx context.Context, chainID *big.Int, call geth.CallMsg) (uint64, error) {
	return m.ec.EstimateGas(ctx, chainID, call)
}
