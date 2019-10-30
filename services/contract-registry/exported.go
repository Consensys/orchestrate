package contractregistry

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry/pg"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry/redis"
)

const (
	component   = "contract-registry"
	postgresOpt = "postgres"
	redisOpt    = "redis"
	mockOpt     = "mock"
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
		case redisOpt:
			// Initialize redis Registry
			redis.Init()

			// Create contract-registry
			registry = redis.GlobalContractRegistry()
		case mockOpt:
			// Initialize mock Registry
			mock.Init()

			// Create contract-registry
			registry = mock.GlobalContractRegistry()
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
