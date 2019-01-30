package ethereum

import (
	"context"
	"math/big"
	"sync/atomic"
	"testing"

	geth "github.com/ethereum/go-ethereum"
)

type MockGasManagerEthClient struct {
	sgpCount uint64
	egCount  uint64
}

func (ec *MockGasManagerEthClient) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	atomic.AddUint64(&ec.sgpCount, 1)
	return nil, nil
}

func (ec *MockGasManagerEthClient) EstimateGas(ctx context.Context, chainID *big.Int, call geth.CallMsg) (uint64, error) {
	atomic.AddUint64(&ec.egCount, 1)
	return 0, nil
}

func TestGasManager(t *testing.T) {
	ec := MockGasManagerEthClient{}

	gm := NewGasManager(&ec)
	gm.SuggestGasPrice(context.Background(), big.NewInt(10))
	if ec.sgpCount != 1 {
		t.Errorf("Expected calls count to SuggestGasPrice to be 1 but got %v", ec.sgpCount)
	}

	gm.EstimateGas(context.Background(), big.NewInt(10), geth.CallMsg{})
	if ec.egCount != 1 {
		t.Errorf("Expected calls count to EstimateGas to be 1 but got %v", ec.egCount)
	}
}
