package infra

import (
	"sync"
	"testing"

	"github.com/Shopify/sarama"
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

func TestSaramaMarshallerConcurrent(t *testing.T) {
	u := SaramaMarshaller{}
	messages := make([]*sarama.ProducerMessage, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		msg := &sarama.ProducerMessage{}
		messages = append(messages, msg)
		wg.Add(1)
		go func(msg *sarama.ProducerMessage) {
			defer wg.Done()
			u.Marshal(newProtoMessage(), msg)
		}(msg)
	}
	wg.Wait()

	for _, msg := range messages {
		b, err := msg.Value.Encode()
		if err != nil {
			t.Errorf("SaramaMarshaller: expected valid value")
		}
		if len(b) < 5 {
			t.Errorf("SaramaMarshaller: expected a non nil message value")
		}
	}
}
