package infra

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore"
)

func initVault(infra *Infra) error {
	ks, err := keystore.NewKeyStore()
	if err != nil {
		log.WithError(err).Fatalf("infra-vault: could not initialize hashicorp Vault")
	}
	log.Infof("infra-vault: vault ready")

	infra.KeyStore = ks

	return nil
}
