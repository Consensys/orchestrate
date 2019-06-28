package sarama

import (
	"sync"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

func newConsumerMessage() *sarama.Msg {
	msg := sarama.Msg{}
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
		if e.GetFrom().Address().Hex() != "0xdbb881a51CD4023E4400CEF3ef73046743f08da3" {
			t.Errorf("Unmarshaller: expected %q but got %q", "abcde", e.From)
		}
	}

}

func TestUnmarshallerError(t *testing.T) {
	msg := &sarama.Msg{
		Value: []byte{0xab, 0x10},
	}
	pb := &ethereum.TxData{}
	e, ok := Unmarshal(msg, pb).(*err.Error)
	assert.NotNil(t, e, "Unmarshal should error")
	assert.True(t, ok, "Error should be internal format")
	assert.Equal(t, e.GetComponent(), "encoding.sarama", "Error code should be correct")
}
