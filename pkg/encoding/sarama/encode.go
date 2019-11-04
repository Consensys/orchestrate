package sarama

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
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
