package engine

import (
	"fmt"
	"sync"
	"testing"
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

func makeMarkerContext(i int) *TxContext {
	txctx := NewTxContext()
	txctx.Reset()
	switch i % 2 {
	case 0:
		txctx.Msg = "error"
		txctx.Keys["errors"] = 1
	case 1:
		txctx.Msg = "valid"
		txctx.Keys["errors"] = 0
	}
	return txctx
}

func TestMarker(t *testing.T) {
	mo := MockOffsetMarker{t: t}
	marker := Marker(&mo)

	rounds := 100
	outs := make(chan *TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		txctx := makeMarkerContext(i)
		go func(txctx *TxContext) {
			defer wg.Done()
			marker(txctx)
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
