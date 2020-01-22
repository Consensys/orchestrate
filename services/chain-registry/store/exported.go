package store

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

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

		// Init Config
		nodes := viper.GetStringSlice(InitViperKey)
		for _, v := range nodes {
			node := &types.Node{
				// Default values
				ListenerDepth:           1,
				ListenerBlockPosition:   -1,
				ListenerFromBlock:       -1,
				ListenerBackOffDuration: "1s",
			}
			dec := json.NewDecoder(strings.NewReader(v))
			dec.DisallowUnknownFields() // Force errors if unknown fields
			err := dec.Decode(node)
			if err != nil {
				log.Warnf("%s: init - invalid node config - got %v", component, err)
			}
			if !viper.GetBool(multitenancy.EnabledViperKey) {
				node.TenantID = multitenancy.DefaultTenantIDName
			}

			err = store.RegisterNode(ctx, node)
			if err != nil {
				log.Warnf("%s: init - could not register node - got %v", component, err)
			}

			log.Infof("%s: init - node %s registered with id %s", component, node.Name, node.ID)
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
