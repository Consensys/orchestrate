package infra

import (
	"sync"
	"testing"

	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

func newTrace() *trace.Trace {
	return &trace.Trace{
		Sender: &common.Account{Id: "abcde"},
	}
}

func TestTracePbMarshaller(t *testing.T) {
	m := TracePbMarshaller{}
	pbs := make([]*trace.Trace, 0)
	rounds := 1000

	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		pbs = append(pbs, &trace.Trace{})
		wg.Add(1)
		go func(pb *trace.Trace) {
			defer wg.Done()
			m.Marshal(newTrace(), pb)
		}(pbs[len(pbs)-1])
	}
	wg.Wait()

	for _, pb := range pbs {
		if pb.Sender.Id != "abcde" {
			t.Errorf("TracePbMarshaller: expected %q but got %q", "abcde", pb.GetSender().GetId())
		}
	}
}
