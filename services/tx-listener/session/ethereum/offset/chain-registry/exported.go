package chainregistry

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
)

var (
	manager  *Manager
	initOnce = &sync.Once{}
)

// Init Offset manager
func Init() {
	initOnce.Do(func() {

		client.Init()
		manager = NewManager(client.GlobalClient())
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
