package infra

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// DummyStore does not store re-cycles new traces
type DummyStore struct {
	pool *sync.Pool
}

// NewDummyStore creates a DummyStore
func NewDummyStore() *DummyStore {
	return &DummyStore{
		pool: &sync.Pool{
			New: func() interface{} { return types.NewTrace() },
		},
	}
}

// Store re-cycle trace
func (s *DummyStore) Store(t *types.Trace) error {
	s.pool.Put(t)
	return nil
}

// Load retrieve a new empty trace
func (s *DummyStore) Load(key interface{}) (*types.Trace, error) {
	t := s.pool.Get().(*types.Trace)
	t.Reset()
	return t, nil
}
