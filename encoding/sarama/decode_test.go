package sarama

import (
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

func newConsumerMessage() *sarama.ConsumerMessage {
	msg := sarama.ConsumerMessage{}
	msg.Value, _ = proto.Marshal(testEnvelope)
	return &msg
}

func TestUnmarshaller(t *testing.T) {
	envelopes := make([]*envelope.Envelope, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		envelopes = append(envelopes, &envelope.Envelope{})
		wg.Add(1)
		go func(e *envelope.Envelope) {
			defer wg.Done()
			_ = Unmarshal(newConsumerMessage(), e)
		}(envelopes[len(envelopes)-1])
	}
	wg.Wait()

	for _, e := range envelopes {
		if e.Sender.Id != "abcde" {
			t.Errorf("Unmarshaller: expected %q but got %q", "abcde", e.Sender.Id)
		}
	}

}
