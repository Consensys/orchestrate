package cache

import (
	"context"
	"time"
)

//go:generate mockgen -source=cache.go -destination=mocks/cache.go -package=mocks

type Manager interface {
	Get(context.Context, string) ([]byte, bool)
	Set(context.Context, string, []byte) bool
	Delete(context.Context, string)
	SetWithTTL(context.Context, string, []byte, time.Duration) bool
	TTL() time.Duration
}
