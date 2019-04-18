package secretstore

import (
	"context"
	"sync"
)

var (
	store    SecretStore
	initOnce = &sync.Once{}
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		// TODO: implement wire by basing on exported in keystore
	})
}

// SetGlobalHandler sets global Faucet Handler
func SetGlobalHandler(s SecretStore) {
	store = s
}

// GlobalHandler returns global Faucet handler
func GlobalHandler() SecretStore {
	return store
}
