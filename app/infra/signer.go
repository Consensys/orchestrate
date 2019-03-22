package infra

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
)

func initSigner(infra *Infra) error {
	// Create Vault Config
	config := secretstore.NewConfig()
	hashicorp, err := secretstore.NewHashicorps(config)
	if err != nil {
		log.WithError(err).Fatalf("infra-vault: could not initialize hashicorp Vault")
	}

	// Initialize hashicorp Vault
	err = secretstore.AutoInit(hashicorp)
	if err != nil {
		log.WithError(err).Fatalf("infra-vault: could not initialize hashicorp Vault")
	}

	infra.SecretStore = hashicorp

	// Declare Secret Store and pre-register private keys
	ks := keystore.NewBaseKeyStore(hashicorp)
	ks.RegisterPkeys(viper.GetStringSlice("secret.pkeys"))
	infra.KeyStore = ks

	log.Infof("infra-vault: ready")

	return nil
}
