package handlers

import (
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

func TestComputeTxCost(t *testing.T) {
	tx := types.NewTx()
	tx.SetGasLimit(32)
	tx.SetGasPrice(big.NewInt(32))

	cost := txCost(tx)

	if cost.Int64() != 1024 {
		t.Errorf("txCost: expected %v but got %v", 1024, cost.Int64())
	}
}

func TestSimpleCreditController(t *testing.T) {
	cfg := &SimpleCreditControllerConfig{
		balanceAt: func(chainID *big.Int, a common.Address) (*big.Int, error) {
			if a.Hex() == "0xdbb881a51CD4023E4400CEF3ef73046743f08da3" {
				return big.NewInt(2000), nil
			}
			return big.NewInt(1000), nil
		},
		creditAmount: big.NewInt(1000),
		maxBalance:   big.NewInt(1999),
		creditDelay:  time.Duration(60 * time.Second),
		blackList:    map[string]struct{}{"0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff": struct{}{}},
	}
	c := NewSimpleCreditController(cfg, 10)

	// Black listed address should not be credited
	amount, ok := c.ShouldCredit(big.NewInt(10), common.HexToAddress("0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"), big.NewInt(100))
	if ok || amount != nil {
		t.Errorf("SimpleCreditController: should not credit black listed account")
	}

	// Black listed address should not be credited
	amount, ok = c.ShouldCredit(big.NewInt(10), common.HexToAddress("0xdbb881a51CD4023E4400CEF3ef73046743f08da3"), big.NewInt(100))
	if ok || amount != nil {
		t.Errorf("SimpleCreditController: should not credit account with to high balance")
	}

	// Should credit in nominal case
	amount, ok = c.ShouldCredit(big.NewInt(10), common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684"), big.NewInt(100))
	if !ok || amount.Int64() != 1000 {
		t.Errorf("SimpleCreditController: should credit in nominal case")
	}

	// Should not credit befor delay
	amount, ok = c.ShouldCredit(big.NewInt(10), common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684"), big.NewInt(100))
	if ok || amount != nil {
		t.Errorf("SimpleCreditController: should not credit account that has not cooldown")
	}
}

type mockEthCrediter struct {
	count int64
	t     *testing.T
}

func (c *mockEthCrediter) Credit(chainID *big.Int, a common.Address, value *big.Int) error {
	c.t.Logf("%v Crediting %q on %v", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), a.Hex(), chainID.Text(16))
	atomic.AddInt64(&c.count, 1)
	return nil
}

var testData = []struct {
	chainID string
	a       string
}{
	{"0x1", "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"},
	{"0x1", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
	{"0x1", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
	{"0x2", "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"},
	{"0x1", "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"},
	{"0x2", "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"},
}

func newFaucetTestMessage(i int) *tracepb.Trace {
	var pb tracepb.Trace
	pb.Chain = &tracepb.Chain{Id: testData[i%6].chainID}
	pb.Sender = &tracepb.Account{Address: testData[i%6].a}
	pb.Transaction = &ethpb.Transaction{
		TxData: &ethpb.TxData{
			Value:    "0xa2bfe3",
			GasPrice: "0xa2b",
			Gas:      1000,
		},
	}
	return &pb
}

func TestFaucet(t *testing.T) {
	// Create worker
	w := infra.NewWorker(100)
	w.Use(Loader(&TraceProtoUnmarshaller{}))

	// Create and register Faucet handler
	cfg := &SimpleCreditControllerConfig{
		balanceAt: func(chainID *big.Int, a common.Address) (*big.Int, error) {
			if a.Hex() == "0xdbb881a51CD4023E4400CEF3ef73046743f08da3" {
				return big.NewInt(2000), nil
			}
			return big.NewInt(1000), nil
		},
		creditAmount: big.NewInt(1000),
		maxBalance:   big.NewInt(1999),
		creditDelay:  time.Duration(10 * time.Millisecond),
		blackList:    map[string]struct{}{"0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff": struct{}{}},
	}
	crediter, controller := &mockEthCrediter{t: t}, NewSimpleCreditController(cfg, 50)
	h := Faucet(crediter, controller)
	w.Use(h)

	// Create input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed input channel and then close it
	// `newFaucetTestMessage(i)` as been designed such as 1 out of 6 message are valid for a credit
	rounds := 600
	for i := 1; i <= rounds; i++ {
		in <- newFaucetTestMessage(i)
		if i%6 == 0 {
			// Sleep to reboot delay on controller
			time.Sleep(10 * time.Millisecond)
		}
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	// Run worker
	go w.Run(in)

	// We ensure that one out of 6 message have been credited
	if crediter.count != 100 {
		t.Errorf("Faucet: expected %v credits but got %v", 100, crediter.count)
	}
}
