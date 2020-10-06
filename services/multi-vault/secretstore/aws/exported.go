package aws

import (
	"context"
	"sync"

	healthz "github.com/heptiolabs/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

const component = "secret-store.aws"

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

		multitenancy.Init(ctx)
		// Set store
		store = NewSecretStore(multitenancy.GlobalKeyBuilder())

		checker = store.Health
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

func GlobalChecker() healthz.Check {
	return checker
}
