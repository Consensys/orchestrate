package aws

import (
	"context"
	"sync"
)

const component = "secret-store.aws"

var (
	store    *SecretStore
	initOnce = &sync.Once{}
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		// Set store
		store = NewSecretStore()
	})
}

// SetGlobalStore sets global mock SecretStore
func SetGlobalStore(s *SecretStore) {
	store = s
}

// GlobalStore returns global mock SecretStore
func GlobalStore() *SecretStore {
	return store
}
