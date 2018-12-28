package handlers

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

type MockNonceLocker struct {
	t      *testing.T
	mux    *sync.Mutex
	value  uint64
	status uint
}

func (l *MockNonceLocker) Lock() error {
	l.mux.Lock()
	if l.status == 2 {
		l.mux.Unlock()
		return fmt.Errorf("Could not lock")
	}
	return nil
}

func (l *MockNonceLocker) Unlock() error {
	defer l.mux.Unlock()
	if l.status == 3 {
		return fmt.Errorf("Could not unlock")
	}
	return nil
}

func (l *MockNonceLocker) Get() (uint64, error) {
	if l.status == 4 {
		return 0, fmt.Errorf("Could not Get")
	}
	return l.value, nil
}

func (l *MockNonceLocker) Set(v uint64) error {
	if l.status == 5 {
		return fmt.Errorf("Could not Set")
	}
	return nil
}

type MockNonceManager struct {
	t       *testing.T
	mux     *sync.Mutex
	counter uint
}

func (m *MockNonceManager) Obtain(chainID *big.Int, a common.Address) (infra.NonceLocker, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.counter%6 == 1 {
		m.counter++
		return nil, fmt.Errorf("Could not obtain")
	}
	n := &MockNonceLocker{
		t:      m.t,
		mux:    m.mux,
		value:  0,
		status: m.counter % 6,
	}
	m.counter++
	return n, nil
}

func makeNonceContext(i int) *infra.Context {
	ctx := infra.NewContext()
	ctx.Reset()
	return ctx
}

func TestNonce(t *testing.T) {
	m := MockNonceManager{t: t, mux: &sync.Mutex{}, counter: 0}
	nonceH := NonceHandler(&m)

	rounds := 100
	outs := make(chan *infra.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeNonceContext(i)
		go func(ctx *infra.Context) {
			defer wg.Done()
			nonceH(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("NonceHandler: expected %v outs but got %v", rounds, len(outs))
	}
}

// import (
// 	"math/big"
// 	"math/rand"
// 	"testing"
// 	"time"

// 	"github.com/ethereum/go-ethereum/common"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
// 	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
// )

// func newNonceTest(chainID *big.Int, a common.Address) (uint64, error) {
// 	return 0, nil
// }

// var (
// 	chains  = []string{"0xa1bc", "0xde4f"}
// 	senders = []string{
// 		"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
// 		"0x0115cB08B395C2c02b82FaD44a698EFA0f47F15f",
// 		"0xb60b036e4fedec7411b6F85E53Bf883BDE23A2c3",
// 	}
// )

// func newNonceTestMessage(i int) *tracepb.Trace {
// 	var pb tracepb.Trace
// 	pb.Chain = &tracepb.Chain{Id: chains[i%2]}
// 	pb.Sender = &tracepb.Account{Address: senders[i%3]}
// 	return &pb
// }

// func dummyTimeHandler(maxtime int) infra.HandlerFunc {
// 	return func(ctx *infra.Context) {
// 		// Simulate some io time
// 		r := rand.Intn(maxtime)
// 		time.Sleep(time.Duration(r) * time.Millisecond)
// 	}
// }

// func TestNonceHandler(t *testing.T) {
// 	// Create handler with 4 stripes lock nonce cache handler
// 	m := NewCacheNonce(newNonceTest, 4)
// 	h := NonceHandler(m)

// 	// Create new worker
// 	w := infra.NewWorker(100)
// 	w.Use(Loader(&TraceProtoUnmarshaller{}))
// 	w.Use(h)
// 	w.Use(NewMockHandler(10).Handler())

// 	// Create input channel
// 	in := make(chan interface{})

// 	// Run worker
// 	go w.Run(in)

// 	// Feed input channel and then close it
// 	rounds := 1000
// 	for i := 1; i <= rounds; i++ {
// 		in <- newNonceTestMessage(i)
// 	}
// 	close(in)

// 	// Wait for worker to be done
// 	<-w.Done()

// 	// Run worker
// 	go w.Run(in)

// 	// Ensure nonces have been properly updated
// 	keys := []struct {
// 		key   string
// 		count uint64
// 	}{
// 		{"a1bc-0xfF778b716FC07D98839f48DdB88D8bE583BEB684", 166},
// 		{"de4f-0x0115cB08B395C2c02b82FaD44a698EFA0f47F15f", 167},
// 		{"a1bc-0xb60b036e4fedec7411b6F85E53Bf883BDE23A2c3", 167},
// 		{"de4f-0xfF778b716FC07D98839f48DdB88D8bE583BEB684", 167},
// 		{"a1bc-0x0115cB08B395C2c02b82FaD44a698EFA0f47F15f", 167},
// 		{"de4f-0xb60b036e4fedec7411b6F85E53Bf883BDE23A2c3", 166},
// 	}
// 	for _, key := range keys {
// 		nonce, ok := m.nonces.Load(key.key)
// 		if !ok {
// 			t.Errorf("NonceHandler: expected nonce on key=%q to have been incremented", key.key)
// 		}
// 		n := nonce.(*SafeNonce)
// 		if n.value != key.count {
// 			t.Errorf("NonceHanlder: expected Nonce %v bug got %v", key.count, n.value)
// 		}
// 	}
// }
