package sarama

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
)

// Unmarshal message assumed to be a sarama.ConsumerMessage into a proto
func Unmarshal(msg *sarama.Msg, pb proto.Message) error {
	// Unmarshal Sarama message to Envelope
	err := proto.Unmarshal(msg.Value, pb)
	if err != nil {
		return err
	}

	return nil
}
