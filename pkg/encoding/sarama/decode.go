package sarama

import (
	"github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	encoding "github.com/ConsenSys/orchestrate/pkg/encoding/proto"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/golang/protobuf/proto"
)

// Unmarshal a sarama message into a protobuffer
func Unmarshal(msg *sarama.Msg, pb proto.Message) error {
	// Unmarshal Sarama message to Envelope
	err := encoding.Unmarshal(msg.Value(), pb)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}
	return nil
}
