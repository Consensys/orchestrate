package nonce

import (
	"sync"
)

type RecoveryTracker struct {
	mux      *sync.Mutex
	counters map[string]uint64
}

func NewRecoveryTracker() *RecoveryTracker {
	return &RecoveryTracker{
		mux:      &sync.Mutex{},
		counters: make(map[string]uint64),
	}
}

func (t *RecoveryTracker) Recovering(key string) (count uint64) {
	t.mux.Lock()
	count = t.counters[key]
	t.mux.Unlock()
	return
}

func (t *RecoveryTracker) Recover(key string) {
	t.mux.Lock()
	t.counters[key]++
	t.mux.Unlock()
}

func (t *RecoveryTracker) Recovered(key string) {
	t.mux.Lock()
	delete(t.counters, key)
	t.mux.Unlock()
}
