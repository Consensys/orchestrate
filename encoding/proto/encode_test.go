package proto

import (
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

func newEnvelope() *envelope.Envelope {
	return &envelope.Envelope{
		Sender: &common.Account{Id: "abcde"},
	}
}

func TestEnvelopeMarshaller(t *testing.T) {
	m := Marshaller{}
	pbs := make([]*envelope.Envelope, 0)
	rounds := 1000

	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		pbs = append(pbs, &envelope.Envelope{})
		wg.Add(1)
		go func(pb *envelope.Envelope) {
			defer wg.Done()
			_ = m.Marshal(newEnvelope(), pb)
		}(pbs[len(pbs)-1])
	}
	wg.Wait()

	for _, pb := range pbs {
		if pb.Sender.Id != "abcde" {
			t.Errorf("EnvelopeMarshaller: expected %q but got %q", "abcde", pb.GetSender().GetId())
		}
	}
}
