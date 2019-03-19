package handlers

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
)

type MockTxSender struct {
	t *testing.T
}

func (s *MockTxSender) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	if chainID.Text(10) == "0" {
		return fmt.Errorf("Could not send")
	}
	return nil
}

func makeSenderContext(i int) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	switch i % 4 {
	case 0:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a")
		ctx.Keys["errors"] = 0
	case 1:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a")
		ctx.Keys["errors"] = 1
	case 2:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw(``)
		ctx.Keys["errors"] = 0
	case 3:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw(``)
		ctx.Keys["errors"] = 0
	}
	return ctx
}

func TestSender(t *testing.T) {
	s := MockTxSender{t: t}
	sender := Sender(&s)

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeSenderContext(i)
		go func(ctx *worker.Context) {
			defer wg.Done()
			sender(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("Marker: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount := out.Keys["errors"].(int)
		if len(out.T.Errors) != errCount {
			t.Errorf("Marker: expected %v errors but got %v", errCount, out.T.Errors)
		}
	}
}
