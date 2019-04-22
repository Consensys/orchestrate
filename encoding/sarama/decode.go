package sarama

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

// Unmarshaller unmarshals Sarama consumer message to an Envelope
type Unmarshaller struct{}

// NewUnmarshaller creates a new marshaller
func NewUnmarshaller() *Unmarshaller {
	return &Unmarshaller{}
}

// Unmarshal message assumed to be a sarama.ConsumerMessage into a proto
func (u *Unmarshaller) Unmarshal(msg *sarama.ConsumerMessage, pb proto.Message) error {
	// Unmarshal Sarama message to Envelope
	err := proto.Unmarshal(msg.Value, pb)
	if err != nil {
		return err
	}

	return nil
}
