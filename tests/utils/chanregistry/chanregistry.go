package chanregistry

import (
	"fmt"
	"sync"

	"github.com/ConsenSys/orchestrate/pkg/types/tx"
)

// ChanRegistry holds a set of indexed envelopes channels
// and allows to dispatch envelope in those channels
type ChanRegistry struct {
	mux *sync.RWMutex

	chans map[string]chan *tx.Envelope
}

// NewChanRegistry creates a new channel registry
func NewChanRegistry() *ChanRegistry {
	return &ChanRegistry{
		mux:   &sync.RWMutex{},
		chans: make(map[string]chan *tx.Envelope),
	}
}

// Register register a new channel
func (r *ChanRegistry) Register(key string, ch chan *tx.Envelope) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.chans[key] = ch
}

// HasChan returns whether a channel is registered for the given key
func (r *ChanRegistry) HasChan(key string) bool {
	r.mux.RLock()
	defer r.mux.RUnlock()

	_, ok := r.chans[key]

	return ok
}

// HasChan returns whether a channel is registered for the given key
func (r *ChanRegistry) GetChan(key string) chan *tx.Envelope {
	r.mux.RLock()
	defer r.mux.RUnlock()

	ch, ok := r.chans[key]
	if !ok {
		return nil
	}

	return ch
}

// Send envelope to channel registered for key
func (r *ChanRegistry) Send(key string, e *tx.Envelope) error {
	r.mux.RLock()
	defer r.mux.RUnlock()

	ch, ok := r.chans[key]
	if !ok {
		return fmt.Errorf("no channel register for key %q", key)
	}

	// Send envelope into channel
	ch <- e

	return nil
}
