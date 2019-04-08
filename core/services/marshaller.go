package services

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
)

// Marshaller are responsible to marshal Envelope object into specific formats (e.g a Sarama message)
type Marshaller interface {
	// Marshal a protobuffer message to specific format
	Marshal(t *envelope.Envelope, msg interface{}) error
}
