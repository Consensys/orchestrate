package handlers

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

type TestGasMsg struct {
	chainID string
	sender  string
	to      string
	value   string
	data    string
}

func newGasTestMessage() *tracepb.Trace {
	var pb tracepb.Trace
	pb.Chain = &tracepb.Chain{Id: "0x1"}
	pb.Sender = &tracepb.Account{Address: "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"}
	pb.Transaction = &ethpb.Transaction{
		TxData: &ethpb.TxData{
			To:    "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
			Value: "0xa2bfe3",
			Data:  "0xabcdef",
		},
	}
	return &pb
}

type DummyGasEstimator struct{}

func (e *DummyGasEstimator) EstimateGas(chainID *big.Int, call ethereum.CallMsg) (uint64, error) {
	r := rand.Intn(50)
	time.Sleep(time.Duration(r) * time.Millisecond)
	return 18, nil
}

type DummyGasPricer struct{}

func (p *DummyGasPricer) SuggestGasPrice(chainID *big.Int) (*big.Int, error) {
	r := rand.Intn(50)
	time.Sleep(time.Duration(r) * time.Millisecond)
	return big.NewInt(123456789), nil
}

func TestGas(t *testing.T) {
	// Create worker
	w := infra.NewWorker(100)
	w.Use(Loader(&TraceProtoUnmarshaller{}))

	// Create and register gas limit handler
	l := GasLimiter(&DummyGasEstimator{})
	w.Use(l)

	// Create and register gas price handler
	p := GasPricer(&DummyGasPricer{})
	w.Use(p)

	mockH := NewMockHandler(50)
	w.Use(mockH.Handler())

	// Create input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed input channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newNonceTestMessage(i)
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	// Run worker
	go w.Run(in)

	if len(mockH.handled) != rounds {
		t.Errorf("Gas: expected %v rounds but got %v", rounds, len(mockH.handled))
	}

	for _, c := range mockH.handled {
		if c.T.Tx().GasLimit() != 18 {
			t.Errorf("Gas: Expected tx GasLimit to be 18 but got %v", c.T.Tx().GasLimit())
		}

		if c.T.Tx().GasPrice().Text(10) != "123456789" {
			t.Errorf("Gas: Expected tx GasPrice to be %q but got %q", "123456789", c.T.Tx().GasPrice().Text(10))
		}
	}
}
