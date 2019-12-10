package memory

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

const component = "chain-registry.store.memory"

var (
	chainRegistry *ChainRegistry
	initOnce      = &sync.Once{}
)

// Initialize Postgres Chain Registry
func Init() {
	initOnce.Do(func() {
		if chainRegistry != nil {
			return
		}

		// Initialize In Memory store
		chainRegistry = NewChainRegistry()
		log.Infof("%s: chain registry ready", component)
	})
}

// GlobalChainRegistry return a chain registry
func GlobalChainRegistry() *ChainRegistry {
	return chainRegistry
}

// SetGlobalChainRegistry sets a new chain registry
func SetGlobalChainRegistry(r *ChainRegistry) {
	chainRegistry = r
}
