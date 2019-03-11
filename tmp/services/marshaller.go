package services

import "gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"

// Marshaller are responsible to marshal a protobuffer message to a higher level message format
type Marshaller interface {
	// Marshal a protobuffer message to a higher level message format
	Marshal(t *types.Trace, msg interface{}) error
}
