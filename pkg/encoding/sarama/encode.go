package sarama

import (
	"github.com/Shopify/sarama"
	encoding "github.com/consensys/orchestrate/pkg/encoding/proto"
	"github.com/consensys/orchestrate/pkg/errors"
	"google.golang.org/protobuf/proto"
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
