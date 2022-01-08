// +build unit

package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/consensys/orchestrate/pkg/toolkit/database/redis"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryTrackerRedis(t *testing.T) {
	mredis, _ := miniredis.Run()
	conf := &redis.Config{
		Expiration: 1,
		Host:       mredis.Host(),
		Port:       mredis.Port(),
	}

	pool, _ := redis.NewPool(conf)
	rt := NewNonceRecoveryTracker(redis.NewClient(pool, conf))

	testKey := "recovery-tracker-redis"
	n := rt.Recovering(testKey)
	assert.Equal(t, uint64(0), n)

	rt.Recover(testKey)
	rt.Recover(testKey)
	n = rt.Recovering(testKey)
	assert.Equal(t, uint64(2), n)

	rt.Recovered(testKey)
	n = rt.Recovering(testKey)
	assert.Equal(t, uint64(0), n)
}
