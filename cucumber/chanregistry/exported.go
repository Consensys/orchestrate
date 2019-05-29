package chanregistry

import (
	"context"
	"sync"
)

var (
	c        *ChanRegistry
	initOnce = &sync.Once{}
)


// Init initialize EnvelopeRegistry
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if c != nil {
			return
		}

		// Create EnvelopeRegistry
		c = NewChanRegistry()

	})
}


// GlobalHandler returns global EnvelopeRegistry
func GlobalChanRegistry() *ChanRegistry {
	return c
}
