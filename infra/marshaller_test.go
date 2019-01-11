package infra

import (
	"sync"
	"testing"

	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
)

func TestTraceProtoMarshallerConcurrent(t *testing.T) {
	u := TraceProtoUnmarshaller{}
	pbs := make([]*tracepb.Trace, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		pb := &tracepb.Trace{}
		pbs = append(pbs, pb)
		wg.Add(1)
		go func(pb *tracepb.Trace) {
			defer wg.Done()
			u.Unmarshal(newProtoMessage(), pb)
		}(pb)
	}
	wg.Wait()

	for _, pb := range pbs {
		if len(pb.GetSender().GetId()) != 5 {
			t.Errorf("TraceProtoMarshaller: expected a 5 long string but got %q", pb.GetSender().GetId())
		}
	}
}
