package services

import (
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// TraceStore is used to store context
type TraceStore interface {
	// Store stores a trace
	// No key is necessary as TraceStore should be able to deduce key from trace object
	Store(t *trace.Trace) error

	// Load should retrieve a trace
	Load(key interface{}) (*trace.Trace, error)
}
