package hashicorp

import (
	"context"
	"sync"

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

		vault, err := NewSecretStore(ConfigFromViper())
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
