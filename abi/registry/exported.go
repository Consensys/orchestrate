package abi

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/abi/registry/static"
)

var (
	registry Registry
	initOnce = &sync.Once{}
)

// Init initialize ABI Registry
func Init() {
	initOnce.Do(func() {
		// Create registry
		registry = static.NewRegistry()

		// Read ABIs from ABI viper configuration
		contracts, err := FromABIConfig()
		if err != nil {
			log.WithError(err).Fatalf("abi: could not initialize ABI registry")
		}

		// Register contracts
		for _, contract := range contracts {
			registry.RegisterContract(contract)
		}
	})
}

// SetGlobalRegistry sets global ABI registry
func SetGlobalRegistry(r Registry) {
	registry = r
}

// GlobalRegistry returns global ABI registry
func GlobalRegistry() Registry {
	return registry
}
