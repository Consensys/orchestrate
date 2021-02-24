// +build unit

package sarama

import (
	"sync"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/types/tx"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	broker "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/errors"
)

func newConsumerMessage() *broker.Msg {
	msg := broker.Msg{}
	msg.ConsumerMessage.Value, _ = proto.Marshal(envlp.TxRequest())
	return &msg
}

func TestUnmarshaller(t *testing.T) {
	envelopes := make([]*tx.TxRequest, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		envelopes = append(envelopes, &tx.TxRequest{})
		wg.Add(1)
		go func(e *tx.TxRequest) {
			defer wg.Done()
			_ = Unmarshal(newConsumerMessage(), e)
		}(envelopes[len(envelopes)-1])
	}
	wg.Wait()

	for _, e := range envelopes {
		if e.GetParams().GetFrom() != "0xdbb881a51CD4023E4400CEF3ef73046743f08da3" {
			t.Errorf("Unmarshaller: expected %q but got %q", "abcde", e.GetId())
		}
	}

}

func TestUnmarshallerError(t *testing.T) {
	msg := &broker.Msg{
		ConsumerMessage: sarama.ConsumerMessage{Value: []byte{0xab, 0x10}},
	}
	pb := &tx.TxRequest{}
	err := errors.FromError(Unmarshal(msg, pb))
	assert.Error(t, err, "Unmarshal should error")
	assert.Equal(t, err.GetComponent(), "encoding.sarama", "Error code should be correct")
}
