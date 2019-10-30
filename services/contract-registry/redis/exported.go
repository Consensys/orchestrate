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

// Init initialize Contract Registry
func Init() {
	initOnce.Do(func() {
		if registry != nil {
			return
		}

		// Initialize gRPC contract-registry
		registry = NewRegistry(NewPool(Config(), Dial))

		log.Infof("%s: store ready", component)
	})
}

// GlobalContractRegistry returns the global contract-registry object
func GlobalContractRegistry() *ContractRegistry {
	return registry
}

// SetGlobalContractRegistry set global contract-registry
func SetGlobalContractRegistry(r *ContractRegistry) {
	registry = r
}
