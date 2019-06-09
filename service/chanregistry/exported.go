package chanregistry

import (
	"context"
	"sync"
)

var (
	chanregistry *EnvelopeChanRegistry
	initOnce     = &sync.Once{}
)

// Init initialize EnvelopeRegistry
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if chanregistry != nil {
			return
		}

		// Create EnvelopeRegistry
		chanregistry = NewChanRegistry()

	})
}

// SetGlobalChanRegistry sets global ChanRegistry
func SetGlobalChanRegistry(c *EnvelopeChanRegistry) {
	chanregistry = c
}

// GlobalHandler returns global EnvelopeRegistry
func GlobalChanRegistry() *EnvelopeChanRegistry {
	return chanregistry
}
