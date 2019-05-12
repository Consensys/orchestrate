package redis

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

// WaitLockReleaseFunc tells what to do when waiting for the lock to be released
type WaitLockReleaseFunc func(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error

// Nonce allows to manage nonce thanks to a redis base
type Nonce struct {
	pool            *redis.Pool
	timeout         int                 // Timeout in millisecond for lock
	waitLockRelease WaitLockReleaseFunc // Add waitLockRelease as struct item instead of method for editing purpose in testing
}

// Creates a new redis pool
func NewPool(port string) *redis.Pool {
	return &redis.Pool{
		// TODO Fine tune those parameters or make them accessible in config file
		MaxIdle:     10000,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", port) },
	}
}

// NewNonce creates a new nonce cache based on redis
func NewNonce(pool *redis.Pool, lockTimeout int) *Nonce {
	return &Nonce{
		pool:            pool,
		waitLockRelease: waitLockRelease,
		timeout:         lockTimeout,
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

// Get returns the value of the nonce from the cache if it exists and returns the last time
// the nonce was gotten or set (idleTime)
// If the nonce does not exist in the cache, the function returns -1 as idleTime
func (nm *Nonce) Get(chainID *big.Int, a *common.Address) (nonce uint64, ok int, err error) {
	idleTime, err := nm.getIdleTime(chainID, a)
	if err != nil {
		return 0, 0, err
	}
	if idleTime != -1 {
		nonce, err := nm.getCache(chainID, a)
		return nonce, idleTime, err
	}
	return 0, idleTime, nil // idleTime == -1, meaning the nonce is not in the cache
}

// Tells if the nonce is already in the cache or not
func (nm *Nonce) getIdleTime(chainID *big.Int, a *common.Address) (int, error) {
	conn := nm.pool.Get()
	defer conn.Close()
	nonceKey := computeKey(chainID, a)
	idleTime, err := redis.Int(conn.Do("OBJECT", "IDLETIME", nonceKey))
	if err != nil {
		if err.Error() == redis.ErrNil.Error() {
			return -1, nil
		}
		return 0, err
	}
	return idleTime, nil
}

// Get the nonce from redis
func (nm *Nonce) getCache(chainID *big.Int, a *common.Address) (uint64, error) {
	conn := nm.pool.Get()
	defer conn.Close()
	nonceKey := computeKey(chainID, a)
	r, err := redis.Uint64(conn.Do("GET", nonceKey))
	if err != nil {
		return 0, err
	}
	return r, nil
}

// Set updates the nonce in the cache with a new value
func (nm *Nonce) Set(chainID *big.Int, a *common.Address, newNonce uint64) error {
	conn := nm.pool.Get()
	defer conn.Close()
	nonceKey := computeKey(chainID, a)
	_, err := conn.Do("SET", nonceKey, newNonce)
	if err != nil {
		return err
	}
	return nil
}

// Lock acquire a lock for a givne chainID and an ethereum address
func (nm *Nonce) Lock(chainID *big.Int, a *common.Address) (string, error) {
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
		return nm.Lock(chainID, a)
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

// Unlock releases a previously acquired lock
func (nm *Nonce) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
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
	return errors.New("lock based on another locking signature, did not unlock the lock")
}
