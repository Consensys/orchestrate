package sarama

import (
	encoding "github.com/ConsenSys/orchestrate/pkg/encoding/proto"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

// Marshal a proto into a e a sarama.ProducerMessage
func Marshal(pb proto.Message, msg *sarama.ProducerMessage) error {
	// Marshal protobuffer into byte
	b, err := encoding.Marshal(pb)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}

	// Set message value
	msg.Value = sarama.ByteEncoder(b)

	return nil
}
