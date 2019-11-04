package pg

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

const component = "contract-registry.pg"

var (
	contractRegistry *ContractRegistry
	initOnce         = &sync.Once{}
)

// Init initialize Postgres Contract Registry
func Init() {
	initOnce.Do(func() {
		if contractRegistry != nil {
			return
		}

		// Initialize Postgres store
		opts := postgres.NewOptions()
		contractRegistry = NewContractRegistryFromPGOptions(opts)
		log.WithFields(log.Fields{
			"db.address":  opts.Addr,
			"db.database": opts.Database,
			"db.user":     opts.User,
		}).Infof("%s: contract registry ready", component)
	})
}

func GlobalContractRegistry() *ContractRegistry {
	return contractRegistry
}

// SetGlobalContractRegistry sets a new contract registry
func SetGlobalContractRegistry(r *ContractRegistry) {
	contractRegistry = r
}
