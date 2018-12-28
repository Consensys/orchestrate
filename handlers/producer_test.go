package handlers

import (
	"fmt"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

type MockTraceProducer struct {
	t *testing.T
}

func (p *MockTraceProducer) Produce(pb *tracepb.Trace) error {
	if pb.GetChain().GetId() == "unknown" {
		return fmt.Errorf("Could not produce")
	}
	return nil
}

func makeProducerContext(i int) *infra.Context {
	ctx := infra.NewContext()
	ctx.Reset()
	switch i % 2 {
	case 0:
		ctx.Pb.Chain = &tracepb.Chain{Id: "unknown"}
		ctx.Keys["errors"] = 1
	case 1:
		ctx.Pb.Chain = &tracepb.Chain{Id: "known"}
		ctx.Keys["errors"] = 0
	}
	return ctx
}

func TestProducer(t *testing.T) {
	mp := MockTraceProducer{t: t}
	producer := Producer(&mp)

	rounds := 100
	outs := make(chan *infra.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeProducerContext(i)
		go func(ctx *infra.Context) {
			defer wg.Done()
			producer(ctx)
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
