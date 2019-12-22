package kafka

import (
	"context"
	"sync"
)

var (
	hook     *Hook
	initOnce = &sync.Once{}
)

// Init Kafka hook
func Init(ctx context.Context) {
	initOnce.Do(func() {
		hook = &Hook{}
	})
}

// SetGlobalHook set global Kafka hook
func SetGlobalHook(hk *Hook) {
	hook = hk
}

// GlobalHook return global Kafka hook
func GlobalHook() *Hook {
	return hook
}
