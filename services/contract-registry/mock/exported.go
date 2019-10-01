package mock

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

const component = "contract-registry.mock"

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
		registry = NewRegistry()

		log.Infof("%q: store ready", component)
	})
}

func GlobalContractRegistry() *ContractRegistry {
	return registry
}

// SetGlobalContractRegistry set global contract-registry
func SetGlobalContractRegistry(r *ContractRegistry) {
	registry = r
}
