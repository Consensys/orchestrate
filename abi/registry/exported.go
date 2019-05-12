package registry

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/static"
)

var (
	rgstr    Registry
	initOnce = &sync.Once{}
)

// Init initialize ABI Registry
func Init() {
	initOnce.Do(func() {
		// Create registry
		rgstr = static.NewRegistry()

		// Read ABIs from ABI viper configuration
		contracts, err := FromABIConfig()
		if err != nil {
			log.WithError(err).Fatalf("abi: could not initialize ABI registry")
		}

		// Register contracts
		for _, contract := range contracts {
			err = rgstr.RegisterContract(contract)

			if err != nil {
				log.WithError(err).Fatalf("abi: could not register ABI")
			}
		}
	})
}

// SetGlobalRegistry sets global ABI registry
func SetGlobalRegistry(r Registry) {
	rgstr = r
}

// GlobalRegistry returns global ABI registry
func GlobalRegistry() Registry {
	return rgstr
}
