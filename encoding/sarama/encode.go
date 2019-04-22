package sarama

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

// Marshaller marshals Envelope to Sarama producer message
type Marshaller struct{}

// NewMarshaller creates a new marshaller
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

// Marshal a proto into a message assume to be a sarama.ProducerMessage
func (m *Marshaller) Marshal(pb proto.Message, msg *sarama.ProducerMessage) error {
	// Marshal protobuffer into byte
	b, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	// Set message value
	msg.Value = sarama.ByteEncoder(b)

	return nil
}
