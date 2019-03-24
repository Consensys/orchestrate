package handlers

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	commonpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

type MockEthCrediter struct {
	count int32
	t     *testing.T
}

func (c *MockEthCrediter) Credit(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
	if r.ChainID.Text(10) == "0" {
		return big.NewInt(0), false, fmt.Errorf("Could not credit")
	}
	atomic.AddInt32(&c.count, 1)
	return nil, false, nil
}

var blackAddress = "0x664895b5fE3ddf049d2Fb508cfA03923859763C6"

func makeFaucetContext(i int) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	switch i % 4 {
	case 0:
		ctx.T.Chain = &commonpb.Chain{}
		ctx.Keys["errors"] = 0
	case 1:
		ctx.T.Chain = (&commonpb.Chain{}).SetID(big.NewInt(0))
		ctx.Keys["errors"] = 1
	case 2:
		ctx.T.Chain = (&commonpb.Chain{}).SetID(big.NewInt(0))
		ctx.T.Sender = &commonpb.Account{Addr: blackAddress}
		ctx.Keys["errors"] = 0
	case 3:
		ctx.T.Chain = (&commonpb.Chain{}).SetID(big.NewInt(1))
		ctx.Keys["errors"] = 0
	}
	return worker.WithContext(context.Background(), ctx)
}

func TestFaucet(t *testing.T) {
	// Create Faucet handler
	mc := &MockEthCrediter{t: t}
	faucet := Faucet(mc, big.NewInt(1000))

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeFaucetContext(i)
		go func(ctx *worker.Context) {
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
