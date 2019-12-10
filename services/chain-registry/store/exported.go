package store

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

const (
	component   = "chain-registry.store"
	postgresOpt = "postgres"
	memoryOpt   = "in-memory"
)

var (
	store    types.ChainRegistryStore
	initOnce = &sync.Once{}
)

// Init initializes a ChainRegistry store
func Init() {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		switch viper.GetString(TypeViperKey) {
		case postgresOpt:
			// Initialize postgres Registry
			pg.Init()

			// Create contract-registry
			store = pg.GlobalChainRegistry()
		case memoryOpt:
			// Initialize mock Registry
			memory.Init()

			// Create contract-registry
			store = memory.GlobalChainRegistry()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(TypeViperKey),
			}).Fatalf("%s: unknown type", component)
		}
	})
}

// SetGlobalRegistry sets global a chain-registry store
func SetGlobalStoreRegistry(r types.ChainRegistryStore) {
	store = r
}

// GlobalRegistry returns global a chain-registry store
func GlobalStoreRegistry() types.ChainRegistryStore {
	return store
}
