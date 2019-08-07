package redis

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gomodule/redigo/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

// WaitLockReleaseFunc tells what to do when waiting for the lock to be released
type WaitLockReleaseFunc func(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error

// Nonce allows to manage nonce thanks to a redis base
type Nonce struct {
	pool            *redis.Pool
	timeout         int                 // Timeout in millisecond for lock
	waitLockRelease WaitLockReleaseFunc // Add waitLockRelease as struct item instead of method for editing purpose in testing
}

// Conn is a wrapper around a redis.Conn that handles internal errors
type Conn struct {
	redis.Conn
}

func Dial(network, address string, options ...redis.DialOption) (redis.Conn, error) {
	conn, err := redis.Dial(network, address, options...)
	if err != nil {
		return conn, errors.ConnectionError(err.Error())
	}
	return Conn{conn}, nil
}

func (conn Conn) Do(commandName string, args ...interface{}) (interface{}, error) {
	reply, err := conn.Conn.Do(commandName, args...)
	if err != nil {
		return reply, errors.ConnectionError(err.Error())
	}
	return reply, nil
}

// Creates a new redis pool
func NewPool(port string) *redis.Pool {
	return &redis.Pool{
		// TODO Fine tune those parameters or make them accessible in config file
		MaxIdle:     10000,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return Dial("tcp", port) },
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

// Get returns the value of the nonce from the cache if it exists
func (nm *Nonce) Get(chainID *big.Int, a *common.Address) (nonce uint64, inCache bool, err error) {
	conn := nm.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("GET", computeKey(chainID, a))
	if err != nil {
		return 0, false, errors.FromError(err).SetComponent(component)
	}

	if reply == nil {
		return 0, false, nil
	}

	r, err := redis.Uint64(reply, nil)
	if err != nil {
		return 0, true, FromRedisError(err).SetComponent(component)
	}
	return r, true, nil
}

// Set updates the nonce in the cache with a new value
func (nm *Nonce) Set(chainID *big.Int, a *common.Address, newNonce uint64) error {
	conn := nm.pool.Get()
	expirationTime := viper.GetInt("redis.nonce.expiration.time")
	defer conn.Close()
	_, err := conn.Do("SETX", computeKey(chainID, a), expirationTime, newNonce)
	if err != nil {
		return errors.FromError(err).SetComponent(component)
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
		return "", errors.FromError(err).SetComponent(component)
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
		return errors.ConnectionError(err.Error()).SetComponent(component)
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
		return errors.FromError(err).SetComponent(component)
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
			return errors.FromError(err).SetComponent(component)
		}
		return nil
	}

	return errors.InternalError("lock keys do not match lock signature").SetComponent(component)
}
