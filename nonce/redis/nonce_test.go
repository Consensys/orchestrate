package redis

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/alicebob/miniredis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

var chainNonce = uint64(42)
var defaultTimeout = 500

func nMock() (nonce *Nonce, addr string, clean func()) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	cleanRedisManager := func() {
		s.Close()
	}

	pool := NewPool(s.Addr())

	return NewNonce(pool, defaultTimeout), s.Addr(), cleanRedisManager
}

func TestGetAndUpdateNonceCache(t *testing.T) {
	cid := big.NewInt(36)
	a := common.HexToAddress("0xabcdabcdabcdabcdabcdabcd")
	// We don't have nonce for this address in cache
	noNonceAddress := common.HexToAddress("0xabcd")
	expirationTime := viper.GetInt("redis.nonce.expiration.time")
	nonceBefore := 42

	mockRedisConn := redigomock.NewConn()
	mockRedisConn.Clear()
	mockRedisConn.Command("SETEX", computeKey(cid, &a), expirationTime, uint64(nonceBefore+1)).Expect("OK")
	mockRedisConn.Command("GET", computeKey(cid, &a)).Expect(int64(nonceBefore + 1))

	mockRedisConn.Command("GET", computeKey(cid, &noNonceAddress)).Expect(nil)

	mockRedisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockRedisConn, nil
		},
		MaxIdle: 10,
	}
	mockWaitLockRelease := func(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error { return nil }
	n := Nonce{pool: mockRedisPool, waitLockRelease: mockWaitLockRelease}

	t.Run("Nonce not in cache", func(t *testing.T) {
		nonce, inCache, err := n.Get(cid, &noNonceAddress)

		assert.NoError(t, err, "Error should be nil")
		assert.False(t, inCache, "Nonce should not be in cache")
		assert.Equal(t, uint64(0), nonce, "Nonce should not be in cache")
	})

	t.Run("Nonce from cache", func(t *testing.T) {
		nonce := uint64(nonceBefore)
		err := n.Set(cid, &a, nonce+1)
		require.Nil(t, err, "Set should not error")

		nonce, inCache, err := n.Get(cid, &a)
		require.Nil(t, err, "Get should not error")
		assert.True(t, inCache, "Nonce should be in cache")
		assert.Equal(t, nonce, chainNonce+1, "Nonce should have been updated")
	})

	t.Run("Test error handling", func(t *testing.T) {
		cid := big.NewInt(36)
		a := common.HexToAddress("0xabcdabcdabcdabcdabcdabcd")
		mockRedisConn := redigomock.NewConn()
		mockRedisConn.Command("SET").ExpectError(fmt.Errorf("test-error"))
		mockRedisConn.Command("SETEX").ExpectError(fmt.Errorf("test-error"))
		mockRedisConn.Command("GET", computeLockName(cid, &a)).Expect("lock-name").ExpectError(fmt.Errorf("test-error"))
		mockRedisConn.Command("GET", computeKey(cid, &a)).ExpectError(fmt.Errorf("test-error"))
		mockRedisConn.Command("DEL", computeLockName(cid, &a)).ExpectError(fmt.Errorf("test-error"))
		mockRedisPool := &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockRedisConn, nil
			},
			MaxIdle: 10,
		}
		mockWaitLockRelease := func(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error { return nil }
		n := Nonce{pool: mockRedisPool, waitLockRelease: mockWaitLockRelease}

		testError := func(err error) {
			e := errors.FromError(err)
			assert.Equal(t, "nonce.redis", e.GetComponent(), "Component should be correct")
			assert.Equal(t, "test-error", e.GetMessage(), "test-error")
		}

		// First EXISTS response
		_, _, err := n.Get(cid, &a)
		testError(err)
		// Second EXISTS response
		_, _, err = n.Get(cid, &a)
		testError(err)
		err = n.Set(cid, &a, uint64(0))
		testError(err)
		_, err = n.Lock(cid, &a)
		testError(err)
		// First GET response
		err = n.Unlock(cid, &a, "lock-name")
		testError(err)
		// Second GET response
		err = n.Unlock(cid, &a, "lock-name")
		testError(err)
	})

	t.Run("Test error handling UpdateCacheNonce and", func(t *testing.T) {
		cid := big.NewInt(36)
		a := common.HexToAddress("0xabcdabcdabcdabcdabcdabcd")
		mockRedisConn := redigomock.NewConn()
		mockRedisPool := &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return mockRedisConn, nil
			},
			MaxIdle: 10,
		}

		mockWaitLockRelease := func(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error { return nil }
		n := Nonce{pool: mockRedisPool, waitLockRelease: mockWaitLockRelease}

		err := n.Set(cid, &a, uint64(0))
		if err == nil {
			t.Error("Error should have been raised")
		}
	})
}

func TestGetLock(t *testing.T) {
	cid := big.NewInt(36)
	a := common.HexToAddress("0xabcdabcdabcdabcdabcdabcd")
	mockRedisConn := redigomock.NewConn()
	mockRedisConn.Command("SET", computeLockName(cid, &a), redigomock.NewAnyData(), "NX", "PX", redigomock.NewAnyData()).Expect("KO").Expect("OK")
	mockRedisConn.GenericCommand("UNSUBSCRIBE").Expect(nil)
	mockRedisConn.GenericCommand("ECHO").Expect(nil)
	mockRedisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockRedisConn, nil
		},
		MaxIdle: 10,
	}
	mockWaitLockRelease := func(chainID *big.Int, a *common.Address, c redis.Conn, timeout time.Duration) error { return nil }
	n := Nonce{pool: mockRedisPool, waitLockRelease: mockWaitLockRelease}

	lockSig, err := n.Lock(cid, &a)
	if err != nil {
		t.Fatalf("Got error %v", err.Error())
	}
	if lockSig == "" {
		t.Fatal("Should have a valid lockSig, got \"\" instead")
	}
}

func TestReleaseLock(t *testing.T) {
	cid := big.NewInt(36)
	a := common.HexToAddress("0xabcdabcdabcdabcdabcdabcd")
	randomIntValue := rand.Int()
	lockSig := strconv.Itoa(randomIntValue)

	n, redisAddr, clean := nMock()
	defer clean()

	// Set the lock
	conn, _ := redis.Dial("tcp", redisAddr)
	_, err := conn.Do("SET", computeLockName(cid, &a), lockSig)
	if err != nil {
		t.Fatalf("Got error %v", err.Error())
	}
	conn.Close()

	// Trying to release the lock with the wrong lockSig
	err = n.Unlock(cid, &a, "wrongLockSig")
	require.NotNil(t, err, "Release should have error on unlocking wrong lock")
	e := errors.FromError(err)
	assert.Equal(t, "nonce.redis", e.GetComponent(), "Component should be correct")
	assert.True(t, errors.IsInternalError(err), "Error should be from correct class")

	err = n.Unlock(cid, &a, lockSig)
	if err != nil {
		t.Fatalf("Got error %v", err.Error())
	}

	// Check what happens if the lock does not exist
	err = n.Unlock(cid, &a, lockSig)
	if err != nil {
		t.Fatalf("Got error %v", err.Error())
	}
}

func TestWaitLockRelease(t *testing.T) {
	cid := big.NewInt(36)
	a := common.HexToAddress("0xabcdabcdabcdabcdabcdabcd")
	c := redigomock.NewConn()
	redisChannel := "__keyspace@*__:" + computeLockName(cid, &a)
	c.Command("PSUBSCRIBE", redisChannel).Expect([]interface{}{
		[]byte("psubscribe"),
		[]byte(redisChannel),
		[]byte("1"),
	})
	c.Command("PUBLISH", redisChannel, "del").Expect([]interface{}{
		[]byte("message"),
		[]byte(redisChannel),
		[]byte("del"),
	})

	timeout := 100 * time.Millisecond

	// Test the normal case
	start := time.Now()
	_ = c.Send("PUBLISH", redisChannel, "del")
	err := waitLockRelease(cid, &a, c, timeout)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	elapsedTime := time.Since(start)
	if elapsedTime > timeout/10 {
		t.Fatal("The function took too long")
	}

	// Test the timeout case
	wrongKey := "wrongKey"
	start = time.Now()
	_ = c.Send("PUBLISH", redisChannel, wrongKey)
	err = waitLockRelease(cid, &a, c, timeout)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	elapsedTime = time.Since(start)
	if elapsedTime < timeout {
		t.Fatal("The function should have timed out")
	}
	if elapsedTime > 2*timeout {
		t.Fatal("The function should not take that long to timeout")
	}
}
