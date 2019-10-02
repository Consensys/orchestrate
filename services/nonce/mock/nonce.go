package mock

import (
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// NonceManager is a NonceManager that works with in memory cache
//
// Important note:
// NonceManager makes the assumption that distinct goroutines access
// nonces for non overlapping set of keys (so there is never competitive access
// to a nonce for a given key)
// Accessing the same key from 2 different goroutines could result
// in discrepancies in nonce updates
type NonceManager struct {
	cache *sync.Map
}

// NewNonceManager creates a new mock NonceManager
func NewNonceManager() *NonceManager {
	return &NonceManager{
		cache: &sync.Map{},
	}
}

const lastAttributedSuf = "last-attributed"

// GetLastAttributed loads last attributed nonce from state
func (nm *NonceManager) GetLastAttributed(key string) (nonce uint64, ok bool, err error) {
	return nm.loadUint64(computeKey(key, lastAttributedSuf))
}

// SetLastAttributed set last attributed nonce
func (nm *NonceManager) SetLastAttributed(key string, value uint64) error {
	nm.set(computeKey(key, lastAttributedSuf), value)
	return nil
}

// IncrLastAttributed increment last attributed nonce
//
// Important note:
// Incrementation does not append atomically so you should
// make sure that IncrLastAttributed is never called by 2 distinct
// goroutines for the same key
func (nm *NonceManager) IncrLastAttributed(key string) (err error) {
	return nm.incrUint64(computeKey(key, lastAttributedSuf))
}

const lastSentSuf = "last-sent"

// GetLastSent loads last sent nonce from state
func (nm *NonceManager) GetLastSent(key string) (nonce uint64, ok bool, err error) {
	return nm.loadUint64(computeKey(key, lastSentSuf))
}

// SetLastSent set last sent nonce
func (nm *NonceManager) SetLastSent(key string, value uint64) error {
	nm.set(computeKey(key, lastSentSuf), value)
	return nil
}

// IncrLastSent increment last sent nonce
//
// Important note:
// Incrementation does not append atomically so you should
// make sure that IncrLastSent is never called by 2 distinct
// goroutines for the same key
func (nm *NonceManager) IncrLastSent(key string) (err error) {
	return nm.incrUint64(computeKey(key, lastSentSuf))
}

const recoveringSuf = "recovering"

// IsRecovering indicates whether NonceManager is in recovery mode
func (nm *NonceManager) IsRecovering(key string) (bool, error) {
	recovering, ok, err := nm.loadBool(computeKey(key, recoveringSuf))
	return recovering && ok, err
}

// SetRecovering set recovery status
func (nm *NonceManager) SetRecovering(key string, status bool) error {
	nm.set(computeKey(key, recoveringSuf), status)
	return nil
}

func (nm *NonceManager) loadUint64(key string) (value uint64, ok bool, err error) {
	v, ok := nm.cache.Load(key)
	if !ok {
		return 0, false, nil
	}

	rv, ok := v.(uint64)
	if !ok {
		return 0, false, errors.InvalidFormatError("loaded value is not uint64")
	}

	return rv, true, nil
}

func (nm *NonceManager) loadBool(key string) (value, ok bool, err error) {
	v, ok := nm.cache.Load(key)
	if !ok {
		return false, false, nil
	}

	rv, ok := v.(bool)
	if !ok {
		return false, false, errors.InvalidFormatError("loaded value is not bool")
	}

	return rv, true, nil
}

func (nm *NonceManager) set(key string, value interface{}) {
	nm.cache.Store(key, value)
}

func (nm *NonceManager) incrUint64(key string) error {
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

func computeKey(key, suffix string) string {
	return fmt.Sprintf("%v-%v", key, suffix)
}
