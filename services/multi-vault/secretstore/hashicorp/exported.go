package hashicorp

import (
	"context"
	"sync"

	healthz "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

const component = "secret-store.hashicorp"

var (
	store    *SecretStore
	initOnce = &sync.Once{}
	checker  healthz.Check
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		// Initialize Key Envelope
		multitenancy.Init(ctx)

		vault, err := NewSecretStore(ConfigFromViper(), multitenancy.GlobalKeyBuilder())
		if err != nil {
			log.Fatalf("KeyStore: Cannot init hashicorp vault got error: %q", err)
		}

		checker = vault.Client.Health

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

func GlobalChecker() healthz.Check {
	return checker
}
