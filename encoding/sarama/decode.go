package sarama

import (
	protobuf "github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// Unmarshal a sarama message into a protobuffer
func Unmarshal(msg *sarama.Msg, pb protobuf.Message) error {
	// Unmarshal Sarama message to Envelope
	err := proto.Unmarshal(msg.Value, pb)
	if err != nil {
		return errors.EncodingError(err.Error()).SetComponent(component)
	}
	return nil
}
