package redis

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

const component = "contract-registry.redis"

var (
	registry *ContractRegistry
	initOnce = &sync.Once{}
)

// Init initialize mock Contract Registry
func Init() {
	initOnce.Do(func() {
		if registry != nil {
			return
		}

		// Initialize gRPC registry
		registry = NewRegistry(NewPool(Config(), Dial))

		log.Infof("%q: store ready", component)
	})
}

// GlobalContractRegistry returns the global contract registry object
func GlobalContractRegistry() *ContractRegistry {
	return registry
}

// SetGlobalContractRegistry set global mock store
func SetGlobalContractRegistry(r *ContractRegistry) {
	registry = r
}
