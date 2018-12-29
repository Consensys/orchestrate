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