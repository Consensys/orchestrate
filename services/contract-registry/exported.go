package contractregistry

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry/redis"
)

var (
	registry svc.RegistryServer
	initOnce = &sync.Once{}
	redisOpt = "redis"
	mockOpt  = "mock"
)

// Init initialize ABI ContractRegistry
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if registry != nil {
			return
		}

		// Initialize Ethereum client
		ethclient.Init(ctx)

		switch viper.GetString(typeViperKey) {
		case redisOpt:
			// Initialize mock Registry
			redis.Init()

			// Create contract-registry
			registry = redis.GlobalContractRegistry()
		case "mock":
			// Initialize mock Registry
			mock.Init()

			// Create contract-registry
			registry = mock.GlobalContractRegistry()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatal("contract-registry: unknown type")
		}

		// Read ABIs from ABI viper configuration
		contracts, err := FromABIConfig()
		if err != nil {
			log.WithError(err).Fatal("abi: could not initialize contract-registry")
		}

		// Register contracts
		for _, contract := range contracts {
			_, err := registry.RegisterContract(ctx, &svc.RegisterContractRequest{Contract: contract})

			if err != nil {
				log.WithError(err).Fatal("abi: could not register ABI")
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
