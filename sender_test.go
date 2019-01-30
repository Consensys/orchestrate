package ethereum

import (
	"context"
	"math/big"
	"sync/atomic"
	"testing"
)

type MockTxSenderEthClient struct {
	srtxCount uint64
}

func (ec *MockTxSenderEthClient) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	atomic.AddUint64(&ec.srtxCount, 1)
	return nil
}

func TestTxSender(t *testing.T) {
	ec := MockTxSenderEthClient{}
	txs := NewTxSender(&ec)
	txs.Send(context.Background(), big.NewInt(10), "0xabcde")
	if ec.srtxCount != 1 {
		t.Errorf("Expected calls count to SendRawTransaction to be 1 but got %v", ec.srtxCount)
	}
}
