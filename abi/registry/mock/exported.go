package mock

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

const component = "contract-registry.mock"

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
		registry = NewRegistry(ethclient.GlobalClient())

		log.Infof("%q: store ready", component)
	})
}

func GlobalContractRegistry() *ContractRegistry {
	return registry
}

// SetGlobalContractRegistry set global mock store
func SetGlobalContractRegistry(r *ContractRegistry) {
	registry = r
}
