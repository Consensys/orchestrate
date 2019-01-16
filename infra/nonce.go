package infra

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gomodule/redigo/redis"
)

const defaultTimeout int = 500

// NonceManager is an interface for fine grain management of nonce by key
type NonceManager interface {
	GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error)
	UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error
	GetLock(chainID *big.Int, a *common.Address) (string, error)
	ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error
}

// RedisWaitLockReleaseFunc tells what to do when waiting for the lock to be released
type RedisWaitLockReleaseFunc func(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error

// RedisNonceManager allows to manage nonce thanks to a redis base
type RedisNonceManager struct {
	pool            *redis.Pool
	timeout         int                      // Timeout in millisecond for lock
	waitLockRelease RedisWaitLockReleaseFunc // Add waitLockRelease as struct item instead of method for editing purpose in testing
}

// Creates a new redis pool
func newRedisPool(port string) *redis.Pool {
	return &redis.Pool{
		// TODO Fine tune those parameters or make them accessible in config file
		MaxIdle:     10000,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", port) },
	}
}

// NewRedisNonceManager creates a new nonce cache based on redis
func NewRedisNonceManager(redisPort string) *RedisNonceManager {
	pool := newRedisPool(redisPort)
	return &RedisNonceManager{
		pool:            pool,
		waitLockRelease: waitLockRelease,
		timeout:         defaultTimeout,
	}
}

// Computes a unique id for a combination of chainID + ethereum address
func computeKey(chainID *big.Int, a *common.Address) string {
	return fmt.Sprintf("%v-%v", chainID.Text(16), a.Hex())
}

// Computes a name for the redis lock given a chainID and an ethereum address
func computeLockName(chainID *big.Int, a *common.Address) string {
	return "lock-" + computeKey(chainID, a)
}

func toDuration(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

// GetNonce returns the value of the nonce (from the cache if available, from the chain otherwise)
func (nm *RedisNonceManager) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	exists, err := nm.nonceInCache(chainID, a)
	if err != nil {
		return 0, false, err
	}
	if exists {
		nonce, err := nm.getCacheNonce(chainID, a)
		return nonce, true, err
	}
	return 0, false, nil
}

// Tells if the nonce is already in the cache or not
func (nm *RedisNonceManager) nonceInCache(chainID *big.Int, a *common.Address) (bool, error) {
	conn := nm.pool.Get()
	defer conn.Close()
	nonceKey := computeKey(chainID, a)
	exists, err := redis.Bool(conn.Do("EXISTS", nonceKey))
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Get the nonce from redis
func (nm *RedisNonceManager) getCacheNonce(chainID *big.Int, a *common.Address) (uint64, error) {
	conn := nm.pool.Get()
	defer conn.Close()
	nonceKey := computeKey(chainID, a)
	r, err := redis.Uint64(conn.Do("GET", nonceKey))
	if err != nil {
		return 0, err
	}
	return r, nil
}

// UpdateCacheNonce updates the nonce in the cache with a new value
func (nm *RedisNonceManager) UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	conn := nm.pool.Get()
	defer conn.Close()
	nonceKey := computeKey(chainID, a)
	_, err := conn.Do("SET", nonceKey, newNonce)
	if err != nil {
		return err
	}
	return nil
}

// GetLock acquire a lock for a givne chainID and an ethereum address
func (nm *RedisNonceManager) GetLock(chainID *big.Int, a *common.Address) (string, error) {
	randomIntValue := rand.Int()
	lockSig := strconv.Itoa(randomIntValue)
	conn := nm.pool.Get()
	defer conn.Close()
	// TODO fix a good value for timeout after tests
	hasLock, err := conn.Do("SET", computeLockName(chainID, a), lockSig, "NX", "PX", strconv.Itoa(nm.timeout))

	if err != nil {
		return "", err
	}

	if hasLock != "OK" {
		err := nm.waitLockRelease(chainID, a, conn, toDuration(nm.timeout))
		conn.Close()
		if err != nil {
			return "", err
		}
		return nm.GetLock(chainID, a)
	}
	return lockSig, nil
}

// wait for a signal saying that the lock has been released or for timeout
func waitLockRelease(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error {
	psc := redis.PubSubConn{Conn: c}
	err := psc.PSubscribe("__keyspace@*__:" + computeLockName(chainID, a))
	if err != nil {
		return err
	}
loop:
	for start := time.Now(); time.Since(start) < timeout; {
		switch n := psc.ReceiveWithTimeout(timeout).(type) {
		case redis.Message:
			if string(n.Data) == "del" {
				break loop
			}
		default:
		}
	}
	return nil
}

// ReleaseLock releases a previously acquired lock
func (nm *RedisNonceManager) ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error {
	conn := nm.pool.Get()
	defer conn.Close()
	lockName := computeLockName(chainID, a)
	raw, err := conn.Do("GET", lockName)
	if err != nil {
		return err
	}
	var value string
	switch raw := raw.(type) {
	case nil:
		// The lock does not exist anymore because of timeout
		return nil
	case []byte:
		value = string(raw)
	case string:
		value = raw
	}
	if value == lockSig {
		_, err := conn.Do("DEL", lockName)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Lock based on another locking signature, did not unlock the lock")
}
