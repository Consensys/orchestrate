package ristretto

import (
	"context"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/cache"
)

type nonCacheManager struct {
}

func NewNonCacheManager() cache.Manager {
	return &nonCacheManager{}
}

func (cca *nonCacheManager) Get(_ context.Context, key string) ([]byte, bool) {
	return nil, false
}

func (cca *nonCacheManager) Set(ctx context.Context, key string, value []byte) bool {
	return true
}

func (cca *nonCacheManager) Delete(_ context.Context, key string) {
}

func (cca *nonCacheManager) SetWithTTL(_ context.Context, key string, value []byte, ttl time.Duration) bool {
	return true
}

func (cca *nonCacheManager) TTL() time.Duration {
	return time.Duration(0)
}
