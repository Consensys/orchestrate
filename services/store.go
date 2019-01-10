package services

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// TraceStore is used to store context
type TraceStore interface {
	// Store stores a trace
	// No key is necessary as TraceStore should be able to deduce key from trace object
	Store(t *types.Trace) error

	// Load should retrieve a trace
	Load(key interface{}) (*types.Trace, error)
}
