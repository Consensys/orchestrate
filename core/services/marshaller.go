package services

import (
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// Marshaller are responsible to marshal trace object into specific formats (e.g a Sarama message)
type Marshaller interface {
	// Marshal a protobuffer message to specific format
	Marshal(t *trace.Trace, msg interface{}) error
}
