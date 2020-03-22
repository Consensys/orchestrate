package memory

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

const component = "secret-store.in-memory"

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

		// Initialize Key Envelope
		multitenancy.Init(ctx)

		// Set store
		store = NewSecretStore(multitenancy.GlobalKeyBuilder())
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
