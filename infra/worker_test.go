package infra

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

type TestHandler struct {
	traces  []*types.Trace
	mux     *sync.Mutex
	handled []*Context
}

func (h *TestHandler) Handler(t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		// We add some randomness in time execution
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}

func newMessage(i uint64) *sarama.ConsumerMessage {
	msg := &sarama.ConsumerMessage{}
	msg.Value, _ = proto.Marshal(
		&tracepb.Trace{
			Transaction: &ethpb.Transaction{
				TxData: &ethpb.TxData{
					Nonce:    i,
					To:       "0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
					Value:    "0x2386f26fc10000",
					Gas:      21136,
					GasPrice: "0xee6b2800",
					Data:     "0xabcd",
				},
			},
		},
	)
	return msg
}

func TestWorker(t *testing.T) {
	h := TestHandler{
		traces:  []*types.Trace{},
		mux:     &sync.Mutex{},
		handled: []*Context{},
	}

	w := NewWorker([]HandlerFunc{h.Handler(t)}, 100)

	// Create a Sarama message channel
	in := make(chan *sarama.ConsumerMessage)

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newMessage(uint64(i))
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if len(h.handled) != rounds {
		t.Errorf("Worker: expected %v rounds but got %v", rounds, len(h.handled))
	}
}
