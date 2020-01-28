package store

import (
	"context"
	"encoding/json"
	"strings"
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
func Init(ctx context.Context) {
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
			return
		}

		// Init chains
		chains := viper.GetStringSlice(InitViperKey)
		importChains(ctx, chains, store)
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

func importChains(ctx context.Context, chains []string, s types.ChainRegistryStore) {
	for _, v := range chains {
		chain := &types.Chain{}
		dec := json.NewDecoder(strings.NewReader(v))
		dec.DisallowUnknownFields() // Force errors if unknown fields
		err := dec.Decode(chain)
		if err != nil {
			log.Warnf("%s: init - invalid chain config - got %v", component, err)
			continue
		}

		err = s.RegisterChain(ctx, chain)
		if err != nil {
			updateErr := s.UpdateChainByName(ctx, chain)
			if updateErr != nil {
				log.Fatalf("%s: init - could not register new chain nor update existing chain - got %v & %v", component, err, updateErr)
			}
			log.Infof("%s: init - chain %s updated", component, chain.Name)
		} else {
			log.Infof("%s: init - chain %s registered with id %s", component, chain.Name, chain.UUID)
		}
	}
}
