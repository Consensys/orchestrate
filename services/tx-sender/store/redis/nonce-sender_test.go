package redis

import (
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/database/redis"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestNonceSenderRedis(t *testing.T) {
	mredis, _ := miniredis.Run()
	conf := &redis.Config{
		Expiration: 1,
		Host:       mredis.Host(),
		Port:       mredis.Port(),
	}

	pool, _ := redis.NewPool(conf)
	ns := NewNonceSender(redis.NewClient(pool, conf))

	testKey := "nonce-sender-redis"
	n, ok, err := ns.GetLastSent(testKey)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), n)

	err = ns.SetLastSent(testKey, 10)
	assert.NoError(t, err)

	n, ok, err = ns.GetLastSent(testKey)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, uint64(10), n)

	err = ns.IncrLastSent(testKey)
	assert.NoError(t, err)
	n, _, _ = ns.GetLastSent(testKey)
	assert.Equal(t, uint64(11), n)

	err = ns.DeleteLastSent(testKey)
	assert.NoError(t, err)
	n, ok, err = ns.GetLastSent(testKey)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), n)
}
