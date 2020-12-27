package memory

import (
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/store"
)

// NonceManager is a NonceManager that works with in memory cache
//
// Important note:
// NonceManager makes the assumption that distinct goroutines access
// nonces for non overlapping set of keys (so there is never competitive access
// to a nonce for a given key)
// Accessing the same key from 2 different goroutines could result
// in discrepancies in nonce updates
type nonceSender struct {
	cache *ristretto.Cache
	ttl   time.Duration
}

// NewNonceSender creates a new mock NonceManager
func NewNonceSender(ttl time.Duration) store.NonceSender {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (100k).
		MaxCost:     1 << 32, // maximum cost of cache (100MB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	return &nonceSender{
		cache: cache,
		ttl:   ttl,
	}
}

const lastSentSuf = "last-sent"

func (nm *nonceSender) GetLastSent(key string) (nonce uint64, ok bool, err error) {
	return nm.loadUint64(computeKey(key, lastSentSuf))
}

// SetLastSent set last sent nonce
func (nm *nonceSender) SetLastSent(key string, value uint64) error {
	nm.set(computeKey(key, lastSentSuf), value)
	return nil
}

func (nm *nonceSender) IncrLastSent(key string) (err error) {
	return nm.incrUint64(computeKey(key, lastSentSuf))
}

func (nm *nonceSender) DeleteLastSent(key string) (err error) {
	nm.delete(computeKey(key, lastSentSuf))
	return nil
}

func (nm *nonceSender) loadUint64(key string) (value uint64, ok bool, err error) {
	v, ok := nm.cache.Get(key)
	if !ok {
		return 0, false, nil
	}

	rv, ok := v.(uint64)
	if !ok {
		return 0, false, errors.InvalidFormatError("loaded value is not uint64")
	}

	return rv, true, nil
}

func (nm *nonceSender) set(key string, value interface{}) {
	size := int64(reflect.TypeOf(value).Size())
	nm.cache.SetWithTTL(key, value, size, nm.ttl)
}

func (nm *nonceSender) delete(key string) {
	nm.cache.Del(key)
}

func (nm *nonceSender) incrUint64(key string) error {
	v, ok, err := nm.loadUint64(key)
	if err != nil {
		return err
	} else if !ok {
		return errors.NotFoundError("no nonce cached for key %q", key)
	}

	// Stores incremented nonce
	nm.set(key, v+1)

	return nil
}
