package handlers

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

type MockNonceGetter struct {
	counter uint64
}

func (g *MockNonceGetter) GetNonce(ctx context.Context, chainID *big.Int, a common.Address) (uint64, error) {
	atomic.AddUint64(&g.counter, 1)
	if chainID.Uint64() == 0 {
		// Simulate error on chain 0
		return 0, fmt.Errorf("Unknwon chain")
	}
	return 42, nil
}

type MockNonceManager struct {
	mux   *sync.Mutex
	nonce *sync.Map
}

func (nm *MockNonceManager) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	if chainID.Uint64() == 1 {
		// Simulate error
		return 0, false, fmt.Errorf("Error retrieving nonce")
	}

	if a.Hex() == "0xfF778b716FC07D98839f48DdB88D8bE583BEB684" {
		// Simulate unknown nonce
		return 0, false, nil
	}

	return 53, true, nil
}

func (nm *MockNonceManager) SetNonce(chainID *big.Int, a *common.Address, value uint64) error {
	if chainID.Uint64() == 2 {
		// Simulate error
		return fmt.Errorf("Error setting nonce")
	}
	return nil
}

func (nm *MockNonceManager) Lock(chainID *big.Int, a *common.Address) (string, error) {
	if chainID.Uint64() == 3 {
		// Simulate error
		return "", fmt.Errorf("Error locking nonce")
	}
	nm.mux.Lock()
	return "random", nil
}

func (nm *MockNonceManager) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
	nm.mux.Unlock()
	if chainID.Uint64() == 4 {
		// Simulate error
		return fmt.Errorf("Error unlocking nonce")
	}
	return nil
}

func makeNonceContext(i int) *types.Context {
	ctx := types.NewContext()
	ctx.Reset()
	ctx.T.Chain().ID = big.NewInt(int64(i % 7))
	if i%7 == 0 || i%7 == 5 {
		*ctx.T.Sender().Address = common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
	} else {
		*ctx.T.Sender().Address = common.HexToAddress("0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff")
	}

	switch i % 7 {
	case 0, 1, 3:
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = uint64(0)
	case 2, 4:
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = uint64(53)
	case 5:
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = uint64(42)
	case 6:
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = uint64(53)
	}

	return ctx
}

func TestNonceHandler(t *testing.T) {
	nm := MockNonceManager{
		mux: &sync.Mutex{},
	}
	ng := MockNonceGetter{}
	nh := NonceHandler(&nm, ng.GetNonce)

	rounds := 100
	outs := make(chan *types.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeNonceContext(i)
		go func(ctx *types.Context) {
			defer wg.Done()
			nh(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("NonceHandler: expected %v outs but got %v", rounds, len(outs))
	}

	for ctx := range outs {
		if len(ctx.T.Errors) != ctx.Keys["errors"].(int) {

			t.Errorf("Expected %v errors but got %v", ctx.Keys["errors"].(int), ctx.T.Errors)
		}
		if ctx.T.Tx().Nonce() != ctx.Keys["result"].(uint64) {
			t.Errorf("Expected Nonce to be %v but got %v", ctx.Keys["result"].(uint64), ctx.T.Tx().Nonce())
		}
	}
}
