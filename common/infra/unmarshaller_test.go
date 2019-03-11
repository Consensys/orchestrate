package infra

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protobuf/trace"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func newProtoMessage() *tracepb.Trace {
	return &tracepb.Trace{
		Sender: &tracepb.Account{Id: "abcde"},
	}
}

func TestTracePbUnmarshaller(t *testing.T) {
	u := TracePbUnmarshaller{}
	traces := make([]*types.Trace, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		traces = append(traces, types.NewTrace())
		wg.Add(1)
		go func(t *types.Trace) {
			defer wg.Done()
			u.Unmarshal(newProtoMessage(), t)
		}(traces[len(traces)-1])
	}
	wg.Wait()

	for _, tr := range traces {
		if tr.Sender().ID != "abcde" {
			t.Errorf("TracePbUnmarshaller: expected %q but got %q", "abcde", tr.Sender().ID)
		}
	}
}
