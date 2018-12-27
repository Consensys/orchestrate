package handlers

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

func newNonceTest(chainID *big.Int, a common.Address) (uint64, error) {
	return 0, nil
}

var (
	chains  = []string{"0xa1bc", "0xde4f"}
	senders = []string{
		"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
		"0x0115cB08B395C2c02b82FaD44a698EFA0f47F15f",
		"0xb60b036e4fedec7411b6F85E53Bf883BDE23A2c3",
	}
)

func newNonceTestMessage(i int) *tracepb.Trace {
	var pb tracepb.Trace
	pb.Chain = &tracepb.Chain{Id: chains[i%2]}
	pb.Sender = &tracepb.Account{Address: senders[i%3]}
	return &pb
}

// func testNonceLoader() infra.HandlerFunc {
// 	return func(ctx *infra.Context) {
// 		msg := ctx.Msg.(*TestNonceMsg)
// 		ctx.Pb.Chain = &tracepb.Chain{Id: msg.chainID}
// 		ctx.Pb.Sender = &tracepb.Account{Address: msg.a}

// 		// Load Trace from protobuffer
// 		protobuf.LoadTrace(ctx.Pb, ctx.T)
// 	}
// }

func dummyTimeHandler(maxtime int) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// Simulate some io time
		r := rand.Intn(maxtime)
		time.Sleep(time.Duration(r) * time.Millisecond)
	}
}

func TestNonceHandler(t *testing.T) {
	// Create handler with 4 stripes lock nonce cache handler
	m := NewCacheNonce(newNonceTest, 4)
	h := NonceHandler(m)

	// Create new worker
	w := infra.NewWorker(100)
	w.Use(TraceProtoLoader())
	w.Use(h)
	w.Use(dummyTimeHandler(10))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newNonceTestMessage(i)
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	// Run worker
	go w.Run(in)

	// Ensure nonces have been properly updated
	keys := []struct {
		key   string
		count uint64
	}{
		{"a1bc-0xfF778b716FC07D98839f48DdB88D8bE583BEB684", 166},
		{"de4f-0x0115cB08B395C2c02b82FaD44a698EFA0f47F15f", 167},
		{"a1bc-0xb60b036e4fedec7411b6F85E53Bf883BDE23A2c3", 167},
		{"de4f-0xfF778b716FC07D98839f48DdB88D8bE583BEB684", 167},
		{"a1bc-0x0115cB08B395C2c02b82FaD44a698EFA0f47F15f", 167},
		{"de4f-0xb60b036e4fedec7411b6F85E53Bf883BDE23A2c3", 166},
	}
	for _, key := range keys {
		nonce, ok := m.nonces.Load(key.key)
		if !ok {
			t.Errorf("NonceHandler: expected nonce on key=%q to have been incremented", key.key)
		}
		n := nonce.(*SafeNonce)
		if n.value != key.count {
			t.Errorf("NonceHanlder: expected Nonce %v bug got %v", key.count, n.value)
		}
	}
}
