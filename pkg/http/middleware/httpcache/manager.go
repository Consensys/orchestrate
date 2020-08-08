package httpcache

import (
	"context"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
)

//go:generate mockgen -source=manager.go -destination=mocks/manager.go -package=mocks
type CacheManager interface {
	Get(context.Context, string) (interface{}, bool)
	Set(context.Context, string, interface{}) bool
	TTL() time.Duration
}

type cacheManager struct {
	c   *ristretto.Cache
	ttl time.Duration
}

func newManager(
	c *ristretto.Cache,
	ttl time.Duration,
) CacheManager {
	return &cacheManager{
		c:   c,
		ttl: ttl,
	}
}

func (cca *cacheManager) Get(_ context.Context, key string) (interface{}, bool) {
	if v, ok := cca.c.Get(key); ok {
		return v, true
	}
	return nil, false
}

func (cca *cacheManager) Set(_ context.Context, key string, value interface{}) bool {
	cost := unsafe.Sizeof(value)
	return cca.c.SetWithTTL(key, value, int64(cost), cca.ttl)
}

func (cca *cacheManager) TTL() time.Duration {
	return cca.ttl
}
