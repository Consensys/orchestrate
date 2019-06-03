package registry

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

var (
	registry Registry
	initOnce = &sync.Once{}
)

// Init initialize ABI Registry
func Init(ctx context.Context) {
	initOnce.Do(func() {
		// Initialize Ethereum client
		ethclient.Init(ctx)

		// Create registry
		registry = static.NewRegistry(ethclient.GlobalClient())

		// Read ABIs from ABI viper configuration
		contracts, err := FromABIConfig()
		if err != nil {
			log.WithError(err).Fatalf("abi: could not initialize ABI registry")
		}

		// Register contracts
		for _, contract := range contracts {
			err = registry.RegisterContract(contract)

			if err != nil {
				log.WithError(err).Fatalf("abi: could not register ABI")
			}
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
