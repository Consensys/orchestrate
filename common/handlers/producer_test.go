package handlers

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
)

type MockProducer struct {
	t *testing.T
}

func (p *MockProducer) Produce(o interface{}) error {
	envelope := o.(*envelope.Envelope)
	if envelope.Chain.ID().Text(10) == "0" {
		return fmt.Errorf("Could not produce")
	}
	return nil
}

func makeProducerContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	switch i % 2 {
	case 0:
		txctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		txctx.Keys["errors"] = 1
	case 1:
		txctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		txctx.Keys["errors"] = 0
	}
	return txctx
}

func TestProducer(t *testing.T) {
	mp := MockProducer{t: t}
	producer := Producer(&mp)

	rounds := 100
	outs := make(chan *engine.TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		txctx := makeProducerContext(i)
		go func(txctx *engine.TxContext) {
			defer wg.Done()
			producer(txctx)
			outs <- txctx
		}(txctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("Marker: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount := out.Keys["errors"].(int)
		if len(out.Envelope.Errors) != errCount {
			t.Errorf("Marker: expected %v errors but got %v", errCount, out.Envelope.Errors)
		}
	}
}
