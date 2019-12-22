package memory

import (
	"context"
	"sync"
)

var (
	manager  *Manager
	initOnce = &sync.Once{}
)

// Init Kafka hook
func Init(ctx context.Context) {
	initOnce.Do(func() {
		manager = NewManager()
	})
}

// SetGlobalHook set global Kafka hook
func SetGlobalManager(mngr *Manager) {
	manager = mngr
}

// GlobalHook return global Kafka hook
func GlobalManager() *Manager {
	return manager
}
