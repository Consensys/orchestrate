package keystore

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore"
)

var (
	keyStore KeyStore
	initOnce = &sync.Once{}
)

// Init initialize Key Store
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if keyStore != nil {
			return
		}

		multitenancy.Init(ctx)
		secretstore.Init(ctx)
		keyStore = base.NewKeyStore(secretstore.GlobalSecretStore())

		err := ImportPrivateKey(keyStore, viper.GetStringSlice(secretPkeyViperKey))
		if err != nil {
			log.Fatalf("Key Store: Cannot import private keys, got error: %q", err)
		}

		log.Info("Key Store: ready")
	})
}

// SetGlobalKeyStore sets global Key Store
func SetGlobalKeyStore(k KeyStore) {
	keyStore = k
}

// GlobalKeyStore returns global Key Store
func GlobalKeyStore() KeyStore {
	return keyStore
}
