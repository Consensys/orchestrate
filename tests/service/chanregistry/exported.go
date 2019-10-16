package chanregistry

import (
	"context"
	"sync"
)

var (
	chanRegistry *ChanRegistry
	initOnce     = &sync.Once{}
)

// Init initialize EnvelopeRegistry
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if chanRegistry != nil {
			return
		}

		// Create EnvelopeRegistry
		chanRegistry = NewChanRegistry()
	})
}

// SetGlobalChanRegistry sets global ChanRegistry
func SetGlobalChanRegistry(c *ChanRegistry) {
	chanRegistry = c
}

// GlobalHandler returns global EnvelopeRegistry
func GlobalChanRegistry() *ChanRegistry {
	return chanRegistry
}
