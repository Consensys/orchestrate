package handlers

import (
	"fmt"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

type MockOffsetMarker struct {
	t *testing.T
}

func (o *MockOffsetMarker) Mark(msg interface{}) error {
	if msg.(string) == "error" {
		return fmt.Errorf("Could not mark")
	}
	return nil
}

func makeMarkerContext(i int) *infra.Context {
	ctx := infra.NewContext()
	ctx.Reset()
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

func TestMarker(t *testing.T) {
	mo := MockOffsetMarker{t: t}
	marker := Marker(&mo)

	rounds := 100
	outs := make(chan *infra.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeMarkerContext(i)
		go func(ctx *infra.Context) {
			defer wg.Done()
			marker(ctx)
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
