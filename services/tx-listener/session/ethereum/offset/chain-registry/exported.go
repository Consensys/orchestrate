package chainregistry

import (
	"context"
	"sync"

	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
)

var (
	manager  *Manager
	initOnce = &sync.Once{}
)

// Init Offset manager
func Init(ctx context.Context) {
	initOnce.Do(func() {

		registry.Init(ctx)
		manager = NewManager(registry.GlobalClient())
	})
}

// SetGlobalHook set global offset manager
func SetGlobalManager(mngr *Manager) {
	manager = mngr
}

// GlobalHook return global offset manager
func GlobalManager() *Manager {
	return manager
}
