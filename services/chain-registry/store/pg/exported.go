package pg

// import (
// 	"sync"

// 	log "github.com/sirupsen/logrus"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
// )

const component = "chain-registry.store.pg"

// var (
// 	chainRegistry *ChainRegistry
// 	initOnce      = &sync.Once{}
// )

// // Initialize Postgres Chain Registry
// func Init() {
// 	initOnce.Do(func() {
// 		if chainRegistry != nil {
// 			return
// 		}

// 		// Initialize Postgres store
// 		opts := postgres.NewOptions()
// 		chainRegistry = NewChainRegistry(postgres.New(opts))
// 		log.WithFields(log.Fields{
// 			"db.address":  opts.Addr,
// 			"db.database": opts.Database,
// 			"db.user":     opts.User,
// 		}).Infof("%s: chain registry ready", component)
// 	})
// }

// // GlobalChainRegistry return a chain registry
// func GlobalChainRegistry() *ChainRegistry {
// 	return chainRegistry
// }

// // SetGlobalChainRegistry sets a new chain registry
// func SetGlobalChainRegistry(r *ChainRegistry) {
// 	chainRegistry = r
// }
