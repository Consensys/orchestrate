package sarama

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

// Marshal a proto into a message assume to be a sarama.ProducerMessage
func Marshal(pb proto.Message, msg *sarama.ProducerMessage) error {
	// Marshal protobuffer into byte
	b, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	// Set message value
	msg.Value = sarama.ByteEncoder(b)

	return nil
}
