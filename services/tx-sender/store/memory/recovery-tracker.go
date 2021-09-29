package memory

import (
	"sync"

	"github.com/consensys/orchestrate/services/tx-sender/store"
)

type nonceRecoveryTracker struct {
	mux      *sync.Mutex
	counters map[string]uint64
}

func NewNonceRecoveryTracker() store.RecoveryTracker {
	return &nonceRecoveryTracker{
		mux:      &sync.Mutex{},
		counters: make(map[string]uint64),
	}
}

const recoverTrackerSuf = "recover-tracker"

func (t *nonceRecoveryTracker) Recovering(key string) (count uint64) {
	t.mux.Lock()
	count = t.counters[computeKey(key, recoverTrackerSuf)]
	t.mux.Unlock()
	return
}

func (t *nonceRecoveryTracker) Recover(key string) {
	t.mux.Lock()
	t.counters[computeKey(key, recoverTrackerSuf)]++
	t.mux.Unlock()
}

func (t *nonceRecoveryTracker) Recovered(key string) {
	t.mux.Lock()
	delete(t.counters, computeKey(key, recoverTrackerSuf))
	t.mux.Unlock()
}
