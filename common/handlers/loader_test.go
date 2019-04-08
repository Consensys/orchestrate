package handlers

import (
	"fmt"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
)

type MockUnmarshaller struct {
	t *testing.T
}

func (u *MockUnmarshaller) Unmarshal(msg interface{}, envelope *envelope.Envelope) error {
	if msg.(string) == "error" {
		return fmt.Errorf("Could not unmarshall")
	}
	return nil
}

func makeLoaderContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Prepare([]engine.HandlerFunc{}, log.NewEntry(log.StandardLogger()), nil)

	switch i % 2 {
	case 0:
		ctx.Msg = "error"
		ctx.Keys["errors"] = 1
	case 1:
		ctx.Msg = "valid"
		ctx.Keys["errors"] = 0
	}
	return ctx
}

func TestLoader(t *testing.T) {
	mu := MockUnmarshaller{t: t}
	loader := Loader(&mu)

	rounds := 10
	outs := make(chan *engine.TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeLoaderContext(i)
		go func(ctx *engine.TxContext) {
			defer wg.Done()
			loader(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)
	if len(outs) != rounds {
		t.Errorf("Loader: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount := out.Keys["errors"].(int)
		if len(out.Envelope.Errors) != errCount {
			t.Errorf("Loader: expected %v errors but got %v", errCount, out.Envelope.Errors)
		}
	}
}
