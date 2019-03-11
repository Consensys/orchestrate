package infra

import (
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protobuf/trace"
)

func newTrace() *types.Trace {
	t := types.NewTrace()
	t.Sender().ID = "abcde"
	return t
}

func TestTracePbMarshaller(t *testing.T) {
	m := TracePbMarshaller{}
	pbs := make([]*tracepb.Trace, 0)
	rounds := 1000

	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		pbs = append(pbs, &tracepb.Trace{})
		wg.Add(1)
		go func(pb *tracepb.Trace) {
			defer wg.Done()
			m.Marshal(newTrace(), pb)
		}(pbs[len(pbs)-1])
	}
	wg.Wait()

	for _, pb := range pbs {
		if pb.GetSender().GetId() != "abcde" {
			t.Errorf("TracePbMarshaller: expected %q but got %q", "abcde", pb.GetSender().GetId())
		}
	}
}
