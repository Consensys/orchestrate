package handlers

import (
	"math/big"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

type DummyTxSender struct {
	mux  *sync.Mutex
	sent []string
}

func (s *DummyTxSender) Send(chainID *big.Int, raw string) error {
	s.mux.Lock()
	s.sent = append(s.sent, raw)
	s.mux.Unlock()
	return nil
}

var testRaw = "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"

func newSenderTestMessage() *tracepb.Trace {
	var pb tracepb.Trace
	pb.Transaction = &ethpb.Transaction{
		Raw: testRaw,
	}
	return &pb
}

func TestSender(t *testing.T) {
	// Create worker
	w := infra.NewWorker(100)

	// Create Sarama loader
	h := Loader(&TraceProtoUnmarshaller{})
	w.Use(h)

	// Register mock handler
	mockH := NewMockHandler(50)
	w.Use(mockH.Handler())

	// Register sender handler
	sender := DummyTxSender{&sync.Mutex{}, []string{}}
	w.Use(Sender(&sender))

	// Create input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed input channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newSenderTestMessage()
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if len(mockH.handled) != rounds {
		t.Errorf("Loader: expected %v rounds but got %v", rounds, len(mockH.handled))
	}

	for _, raw := range sender.sent {
		if raw != testRaw {
			t.Errorf("Loader: expected %q got %q", testRaw, raw)
		}
	}
}
