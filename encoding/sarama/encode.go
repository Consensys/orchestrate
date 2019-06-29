package sarama

import (
	"github.com/Shopify/sarama"
	protobuf "github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/proto"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// Marshal a proto into a e a sarama.ProducerMessage
func Marshal(pb protobuf.Message, msg *sarama.ProducerMessage) error {
	// Marshal protobuffer into byte
	b, err := proto.Marshal(pb)
	if err != nil {
		return errors.EncodingError(err.Error()).SetComponent(component)
	}

	// Set message value
	msg.Value = sarama.ByteEncoder(b)

	return nil
}
