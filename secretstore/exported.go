package secretstore

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/aws"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/hashicorp"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/services"
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
		case "test":
			// Create Key Store from a Mock SecretStore
			mock.Init(ctx)
			secretStore = mock.GlobalStore()

		case "hashicorp":
			// Create an hashicorp vault object
			hashicorp.Init(ctx)
			secretStore = hashicorp.GlobalStore()

		case "aws":
			// Create an hashicorp vault object
			aws.Init(ctx)
			secretStore = aws.GlobalStore()

		default:
			// Key Store type should be one of "test", "hashicorp"
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
