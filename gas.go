package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

// EthGasManager implements methods to get information about Gas by connecting to an Ethereum client
type EthGasManager struct {
	ec *EthClient
}

// NewEthGasManager creates a new EthGasManager
func NewEthGasManager(ec *EthClient) *EthGasManager {
	return &EthGasManager{ec}
}

// SuggestGasPrice suggests a gas price
func (m *EthGasManager) SuggestGasPrice(chainID *big.Int) (*big.Int, error) {
	return m.ec.SuggestGasPrice(context.Background())
}

// EstimateGas suggests a gas limit
func (m *EthGasManager) EstimateGas(chainID *big.Int, call ethereum.CallMsg) (uint64, error) {
	return m.ec.EstimateGas(context.Background(), call)
}
