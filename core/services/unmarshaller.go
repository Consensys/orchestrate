package services

import "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"

// Unmarshaller are responsible to unmarshal high level input message into a protobuf message
type Unmarshaller interface {
	// Unmarshal high level input message into a protobuf message
	Unmarshal(msg interface{}, t *types.Trace) error
}
