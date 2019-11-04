package memory

import (
	"context"
	"sync"
)

var (
	store    *SecretStore
	initOnce = &sync.Once{}
)

// Init initialize Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		// Set store
		store = NewSecretStore()
	})
}

// SetGlobalStore sets global memory SecretStore
func SetGlobalStore(s *SecretStore) {
	store = s
}

// GlobalStore returns global memory SecretStore
func GlobalStore() *SecretStore {
	return store
}
