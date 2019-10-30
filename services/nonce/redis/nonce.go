package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// NonceManager manages nonce using an underlying redis cache
type NonceManager struct {
	pool *redis.Pool
}

// NewNonceManager creates a new NonceManager using an underlying redis cache
func NewNonceManager(pool *redis.Pool) *NonceManager {
	return &NonceManager{
		pool: pool,
	}
}

const lastAttributedSuf = "last-attributed"

// GetLastAttributed loads last attributed nonce from state
func (nm *NonceManager) GetLastAttributed(key string) (nonce uint64, ok bool, err error) {
	return nm.loadUint64(computeKey(key, lastAttributedSuf))
}

// SetLastAttributed set last attributed nonce
func (nm *NonceManager) SetLastAttributed(key string, value uint64) error {
	return nm.set(computeKey(key, lastAttributedSuf), value)
}

// IncrLastAttributed increment last attributed nonce
func (nm *NonceManager) IncrLastAttributed(key string) error {
	return nm.incr(computeKey(key, lastAttributedSuf))
}

const lastSentSuf = "last-sent"

// GetLastSent loads last sent nonce from state
func (nm *NonceManager) GetLastSent(key string) (nonce uint64, ok bool, err error) {
	return nm.loadUint64(computeKey(key, lastSentSuf))
}

// SetLastSent set last set nonce
func (nm *NonceManager) SetLastSent(key string, value uint64) error {
	return nm.set(computeKey(key, lastSentSuf), value)
}

// IncrLastSent increment last sent nonce
func (nm *NonceManager) IncrLastSent(key string) error {
	return nm.incr(computeKey(key, lastSentSuf))
}

const recoveringSuf = "recovering"

// IsRecovering indicates whether NonceManager is in recovery mode
func (nm *NonceManager) IsRecovering(key string) (bool, error) {
	recovering, ok, err := nm.loadBool(computeKey(key, recoveringSuf))
	return recovering && ok, err
}

// SetRecovering set recovery status
func (nm *NonceManager) SetRecovering(key string, status bool) error {
	return nm.set(computeKey(key, recoveringSuf), status)
}

func (nm *NonceManager) load(key string) (value interface{}, ok bool, err error) {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	reply, err := conn.Do("GET", key)
	if err != nil {
		return reply, false, errors.FromError(err).SetComponent(component)
	}

	if reply == nil {
		return nil, false, nil
	}

	return reply, true, nil
}

func (nm *NonceManager) loadUint64(key string) (value uint64, ok bool, err error) {
	// Load value
	reply, ok, err := nm.load(key)
	if err != nil || !ok {
		return 0, false, err
	}

	// Format reply to Uint64
	value, err = redis.Uint64(reply, nil)
	if err != nil {
		return 0, false, FromRedisError(err).SetComponent(component)
	}

	return value, true, nil
}

func (nm *NonceManager) loadBool(key string) (value, ok bool, err error) {
	// Load value
	reply, ok, err := nm.load(key)
	if err != nil || !ok {
		return false, false, err
	}

	// Format reply to Uint64
	value, err = redis.Bool(reply, nil)
	if err != nil {
		return false, false, FromRedisError(err).SetComponent(component)
	}

	return value, true, nil
}

func (nm *NonceManager) set(key string, value interface{}) error {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}

	return nil
}

func (nm *NonceManager) incr(key string) error {
	conn := nm.pool.Get()
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.WithError(closeErr).Warn("could not close redis connection")
		}
	}()

	_, err := conn.Do("INCR", key)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
	}

	return nil
}

func computeKey(key, suffix string) string {
	return fmt.Sprintf("%v-%v", key, suffix)
}
