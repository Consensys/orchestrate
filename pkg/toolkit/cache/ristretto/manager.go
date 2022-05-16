package ristretto

import (
	"context"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/cache"
	"github.com/dgraph-io/ristretto"
)

type cacheManager struct {
	c      *ristretto.Cache
	ttl    *time.Duration
	logger *log.Logger
}

func NewCacheManager(componentID string,
	c *ristretto.Cache,
	ttl *time.Duration,
) cache.Manager {
	return &cacheManager{
		c:      c,
		ttl:    ttl,
		logger: log.NewLogger().SetComponent(componentID),
	}
}

func (cca *cacheManager) Get(_ context.Context, key string) ([]byte, bool) {
	if v, ok := cca.c.Get(key); ok {
		cca.logger.WithField("key", key).Trace("get value")
		return v.([]byte), true
	}
	return nil, false
}

func (cca *cacheManager) Set(ctx context.Context, key string, value []byte) bool {
	if cca.ttl != nil {
		return cca.SetWithTTL(ctx, key, value, *cca.ttl)
	}

	cca.logger.WithField("key", key).Trace("set value")
	return cca.c.Set(key, value, int64(len(value)))
}

func (cca *cacheManager) Delete(_ context.Context, key string) {
	cca.c.Del(key)
}

func (cca *cacheManager) SetWithTTL(_ context.Context, key string, value []byte, ttl time.Duration) bool {
	cca.logger.WithField("key", key).WithField("ttl", ttl.String()).Trace("set value")
	return cca.c.SetWithTTL(key, value, int64(len(value)), ttl)
}

func (cca *cacheManager) TTL() time.Duration {
	if cca.ttl == nil {
		return time.Duration(0)
	}
	return *cca.ttl
}
