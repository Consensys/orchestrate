package hashicorp

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	log "github.com/sirupsen/logrus"
)

const component = "secret-store.hashicorp"

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

		// Initialize Key Builder
		multitenancy.Init(ctx)

		vault, err := NewSecretStore(ConfigFromViper(), multitenancy.GlobalKeyBuilder())
		if err != nil {
			log.Fatalf("Key Store: Cannot init hashicorp vault got error: %q", err)
		}

		// Set store
		store = vault
	})
}

// SetGlobalStore sets global mock SecretStore
func SetGlobalStore(h *SecretStore) {
	store = h
}

// GlobalStore returns global mock SecretStore
func GlobalStore() *SecretStore {
	return store
}
