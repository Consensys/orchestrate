package keystore

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/hashicorp"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/mock"
)

var (
	keyStore  KeyStore
	initOnce = &sync.Once{}
)

// Init initialize Key Store
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if keyStore != nil {
			return
		}

		switch viper.GetString(secretStoreViperKey) {
		case "test":
			// Create Key Store from a Mock SecretStore
			mock.Init(ctx)
			keyStore = base.NewKeyStore(mock.GlobalStore())

		case "hashicorp":
			// Create an hashicorp vault object
			hashicorp.Init(ctx)
			keyStore = base.NewKeyStore(hashicorp.GlobalStore())

		default:
			// Key Store type should be one of "test", "hashicorp"
			log.Fatalf("Key Store: Invalid Store type %q", viper.GetString(secretStoreViperKey))
		}

		err := ImportPrivateKey(keyStore)
		if err != nil {
			log.Fatalf("Key Store: Cannot import private keys, got error: %q", err)
		}

		log.Infof("Key Store: %q ready", viper.GetString(secretStoreViperKey))
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

// ImportPrivateKey create new Key Store
func ImportPrivateKey(k KeyStore) error {

	// Pre-Import Pkeys
	for _, pkey := range viper.GetStringSlice(secretPkeyViperKey) {
		err := k.ImportPrivateKey(pkey)
		if err != nil {
			return err
		}
	}

	return nil
}