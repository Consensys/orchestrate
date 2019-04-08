package services

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
)

// Unmarshaller are responsible to unmarshal input message into an envelope
type Unmarshaller interface {
	// Unmarshal high message into a Envelope
	Unmarshal(msg interface{}, t *envelope.Envelope) error
}
