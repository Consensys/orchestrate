package services

import (
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// Unmarshaller are responsible to unmarshal input message into a trace
type Unmarshaller interface {
	// Unmarshal high message into a trace
	Unmarshal(msg interface{}, t *trace.Trace) error
}
