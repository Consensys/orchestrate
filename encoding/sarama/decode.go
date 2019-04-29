package sarama

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

// Unmarshal message assumed to be a sarama.ConsumerMessage into a proto
func Unmarshal(msg *sarama.ConsumerMessage, pb proto.Message) error {
	// Unmarshal Sarama message to Envelope
	err := proto.Unmarshal(msg.Value, pb)
	if err != nil {
		return err
	}

	return nil
}
