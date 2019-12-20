package aws

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
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

		multitenancy.Init(ctx)
		// Set store
		store = NewSecretStore(multitenancy.GlobalKeyBuilder())
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
