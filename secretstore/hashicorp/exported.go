package hashicorp

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

const component = "secretstore.hashicorp"

var (
	store    *HashiCorp
	initOnce = &sync.Once{}
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		vault, err := NewHashiCorp(NewConfig())
		if err != nil {
			log.Fatalf("Key Store: Cannot init hashicorp vault got error: %q", err)
		}

		// Set store
		store = vault
	})
}

// SetGlobalStore sets global mock SecretStore
func SetGlobalStore(h *HashiCorp) {
	store = h
}

// GlobalStore returns global mock SecretStore
func GlobalStore() *HashiCorp {
	return store
}
