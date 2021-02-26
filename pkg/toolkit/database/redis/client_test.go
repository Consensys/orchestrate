// +build unit

package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func NewRedisMock() *miniredis.Miniredis {
	mredis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return mredis
}

var testKey = "test-key"

func TestRedisClient(t *testing.T) {
	mredis := NewRedisMock()
	conf := &Config{
		Expiration: 1,
		Host:       mredis.Host(),
		Port:       mredis.Port(),
	}

	pool, _ := NewPool(conf)
	nm = NewClient(pool, conf)

	n, ok, err := nm.LoadUint64(testKey)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), n)

	err = nm.Set(testKey, 10)
	assert.NoError(t, err)

	n, ok, err = nm.LoadUint64(testKey)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, uint64(10), n)

	err = nm.Incr(testKey)
	assert.NoError(t, err)
	n, _, _ = nm.LoadUint64(testKey)
	assert.Equal(t, uint64(11), n)
	
	err = nm.Delete(testKey)
	assert.NoError(t, err)
	n, ok, err = nm.LoadUint64(testKey)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), n)
}
