package handlers

import (
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

type MockEthCrediter struct {
	count int32
	t     *testing.T
}

func (c *MockEthCrediter) Credit(chainID *big.Int, a common.Address, value *big.Int) error {
	if chainID.Text(10) == "0" {
		return fmt.Errorf("Could not credit")
	}
	atomic.AddInt32(&c.count, 1)
	return nil
}

type MockEthCreditController struct {
	t *testing.T
}

var blackAddress = "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"

func (c *MockEthCreditController) ShouldCredit(chainID *big.Int, a common.Address, value *big.Int) (*big.Int, bool) {
	if a.Hex() == blackAddress {
		return nil, false
	}
	return big.NewInt(100), true
}

func makeFaucetContext(i int) *infra.Context {
	ctx := infra.NewContext()
	ctx.Reset()
	switch i % 4 {
	case 0:
		ctx.Keys["errors"] = 0
	case 1:
		ctx.T.Chain().ID = big.NewInt(0)
		ctx.Keys["errors"] = 1
	case 2:
		ctx.T.Chain().ID = big.NewInt(0)
		*ctx.T.Sender().Address = common.HexToAddress(blackAddress)
		ctx.Keys["errors"] = 0
	case 3:
		ctx.T.Chain().ID = big.NewInt(1)
		ctx.Keys["errors"] = 0
	}
	return ctx
}

func TestFaucet(t *testing.T) {
	// Create Faucet handler
	mc := &MockEthCrediter{t: t}
	faucet := Faucet(mc, &MockEthCreditController{t: t})

	rounds := 100
	outs := make(chan *infra.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeFaucetContext(i)
		go func(ctx *infra.Context) {
			defer wg.Done()
			faucet(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("Faucet: expected %v outs but got %v", rounds, len(outs))
	}

	if mc.count != int32(rounds/4) {
		t.Errorf("Faucet: expected credit count to be %v but got %v", rounds/4, mc.count)
	}
}
