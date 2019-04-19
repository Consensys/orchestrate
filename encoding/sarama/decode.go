package sarama

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// Unmarshaller unmarshals Sarama consumer message to an Envelope
type Unmarshaller struct{}

// NewUnmarshaller creates a new marshaller
func NewUnmarshaller() *Unmarshaller {
	return &Unmarshaller{}
}

// Unmarshal message
func (u *Unmarshaller) Unmarshal(msg interface{}, e *envelope.Envelope) error {
	// Cast message into a sarama.ConsumerMessage
	cast, ok := msg.(*sarama.ConsumerMessage)
	if !ok {
		return fmt.Errorf("Expected a sarama.ConsumerMessage")
	}

	// Unmarshal Sarama message to Envelope
	err := proto.Unmarshal(cast.Value, e)
	if err != nil {
		return err
	}

	return nil
}
