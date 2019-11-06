package contractregistry

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/pg"
)

const (
	component   = "contract-registry"
	postgresOpt = "postgres"
	memoryOpt   = "in-memory"
)

var (
	registry svc.RegistryServer
	initOnce = &sync.Once{}
)

// Init initialize ABI ContractRegistry
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if registry != nil {
			return
		}

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
func SetGlobalRegistry(r svc.RegistryServer) {
	registry = r
}

// GlobalRegistry returns global contract-registry
func GlobalRegistry() svc.RegistryServer {
	return registry
}
