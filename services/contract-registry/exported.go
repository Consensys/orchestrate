package contractregistry

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/pg"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

const (
	component   = "contract-registry"
	postgresOpt = "postgres"
	memoryOpt   = "in-memory"
)

var (
	registry svc.ContractRegistryServer
	initOnce = &sync.Once{}
)

// Init initialize ABI ContractRegistry
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if registry != nil {
			return
		}

		multitenancy.Init(ctx)

		switch viper.GetString(typeViperKey) {
		case postgresOpt:
			// Initialize postgres Registry
			pg.Init()

			// Create contract-registry
			registry = pg.GlobalContractRegistry()
		case memoryOpt:
			// Initialize mock Registry
			memory.Init()

			// Create contract-registry
			registry = memory.GlobalContractRegistry()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatalf("%s: unknown type", component)
		}

		// Read ABIs from ABI viper configuration
		contracts, err := FromABIConfig()
		if err != nil {
			log.WithError(err).Fatalf("%s: could not initialize contract-registry", component)
		}

		// Register contracts
		for _, contract := range contracts {
			_, err := registry.RegisterContract(ctx, &svc.RegisterContractRequest{Contract: contract})

			if err != nil {
				log.WithError(err).Fatalf("%s: could not register ABI", component)
			}
		}
	})
}

// SetGlobalRegistry sets global contract-registry
func SetGlobalRegistry(r svc.ContractRegistryServer) {
	registry = r
}

// GlobalRegistry returns global contract-registry
func GlobalRegistry() svc.ContractRegistryServer {
	return registry
}
