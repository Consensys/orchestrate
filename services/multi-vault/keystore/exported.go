package keystore

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore"
)

const component = "multi-vault.keystore"

var (
	keyStore keystore.KeyStore
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
		keyStore = NewKeyStore(secretstore.GlobalSecretStore())

		err := importPrivateKey(keyStore, viper.GetStringSlice(secretPkeyViperKey))
		if err != nil {
			log.Fatalf("Key Store: Cannot import private keys, got error: %q", err)
		}

		log.Info("Key Store: ready")
	})
}

// SetGlobalKeyStore sets global Key Store
func SetGlobalKeyStore(k keystore.KeyStore) {
	keyStore = k
}

// GlobalKeyStore returns global Key Store
func GlobalKeyStore() keystore.KeyStore {
	return keyStore
}

// importPrivateKey create new Key Store
func importPrivateKey(k keystore.KeyStore, pkeys []string) error {
	// Pre-Import Pkeys
	for _, pkey := range pkeys {
		ctx, key, err := multitenancy.SplitTenant(pkey)
		if err != nil {
			return err
		}
		err = k.ImportPrivateKey(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}
