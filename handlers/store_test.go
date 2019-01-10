package handlers

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

type MockStore struct {
	mux    *sync.Mutex
	stored []*types.Trace
}

func (s *MockStore) Store(t *types.Trace) error {
	if t.Chain().ID.Text(10) == "0" {
		return fmt.Errorf("Could not store")
	}
	s.mux.Lock()
	s.stored = append(s.stored, t)
	s.mux.Unlock()
	return nil
}

func (s *MockStore) Load(key interface{}) (*types.Trace, error) {
	return s.stored[0], nil
}

func MakeStoreContext(i int) *types.Context {
	ctx := types.NewContext()
	ctx.Reset()
	switch i % 2 {
	case 0:
		ctx.T.Chain().ID = big.NewInt(0)
		ctx.Keys["errors"] = 1
	case 1:
		ctx.T.Chain().ID = big.NewInt(10)
		ctx.Keys["errors"] = 0
	}
	return ctx
}

func TestStore(t *testing.T) {
	// Register sender handler
	ms := MockStore{&sync.Mutex{}, []*types.Trace{}}
	store := Store(&ms)
	rounds := 100
	outs := make(chan *types.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := MakeStoreContext(i)
		go func(ctx *types.Context) {
			defer wg.Done()
			store(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("Store: expected %v outs but got %v", rounds, len(outs))
	}

	if len(ms.stored) != rounds/2 {
		t.Errorf("Store: expected %v stored but got %v", rounds/2, len(ms.stored))
	}
}
