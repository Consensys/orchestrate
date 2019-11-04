package secretstore

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/aws"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/hashicorp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

const (
	hashicorpOpt = "hashicorp"
	memoryOpt    = "in-memory"
	awsOpt       = "aws"
)

var (
	secretStore services.SecretStore
	initOnce    = &sync.Once{}
)

// Init initializes a Secret Store
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if secretStore != nil {
			return
		}

		switch viper.GetString(secretStoreViperKey) {
		case memoryOpt:
			// Create Key Store from a Memory SecretStore
			memory.Init(ctx)
			secretStore = memory.GlobalStore()

		case hashicorpOpt:
			// Create an HashiCorp Vault object
			hashicorp.Init(ctx)
			secretStore = hashicorp.GlobalStore()

		case awsOpt:
			// Create an HashiCorp Vault vault object
			aws.Init(ctx)
			secretStore = aws.GlobalStore()

		default:
			// Key Store type should be one of "memory", "hashicorp"
			log.Fatalf("Secret Store: Invalid Store type %q", viper.GetString(secretStoreViperKey))
		}

		log.Infof("Secret Store: %q ready", viper.GetString(secretStoreViperKey))

	})
}

// SetGlobalHandler sets global Faucet Handler
func SetGlobalSecretStore(s services.SecretStore) {
	secretStore = s
}

// GlobalHandler returns global Faucet handler
func GlobalSecretStore() services.SecretStore {
	return secretStore
}
