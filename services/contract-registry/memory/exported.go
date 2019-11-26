package memory

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

const component = "contract-registry.in-memory"

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
		registry = NewContractRegistry()

		log.Infof("%s: store ready", component)
	})
}

func GlobalContractRegistry() *ContractRegistry {
	return registry
}

// SetGlobalContractRegistry set global contract-registry
func SetGlobalContractRegistry(r *ContractRegistry) {
	registry = r
}
