package chanregistry

import (
	"context"
	"sync"
)

var (
	chanRegistry *EnvelopeChanRegistry
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
func SetGlobalChanRegistry(c *EnvelopeChanRegistry) {
	chanRegistry = c
}

// GlobalHandler returns global EnvelopeRegistry
func GlobalChanRegistry() *EnvelopeChanRegistry {
	return chanRegistry
}
