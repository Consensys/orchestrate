package sarama

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// Marshaller marshals Envelope to Sarama producer message
type Marshaller struct{}

// NewMarshaller creates a new marshaller
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

// Marshal Envelope into a Sarama producer message
func (m *Marshaller) Marshal(e *envelope.Envelope, msg interface{}) error {
	// Cast message into a sarama.ConsumerMessage
	cast, ok := msg.(*sarama.ProducerMessage)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Marshal protobuffer into byte
	b, err := proto.Marshal(e)
	if err != nil {
		return err
	}

	// Set message value
	cast.Value = sarama.ByteEncoder(b)

	return nil
}

