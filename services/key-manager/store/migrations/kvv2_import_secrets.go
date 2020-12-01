package migrations

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	kvv2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/hashicorp/kv-v2"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

func Kvv2ImportSecrets(_ context.Context, vault store.Vault, v2Client *kvv2.Client) error {
	log.Infof("Importing Hashicorp kv-v2 secrets to Vault...")

	if err := vault.HealthCheck(); err != nil {
		return err
	}

	if _, err := v2Client.Health(); err != nil {
		return err
	}

	// Fetch v2 addresses
	addresses, err := v2Client.List("")
	if err != nil {
		// TODO: Check engine is not available and IGNORE if so
		log.WithError(err).Error("could not connect to engine kv-v2")
		return err
	}

	log.Infof("importing accounts %q", addresses)
	for _, addrKey := range addresses {
		privKey, ok, err := v2Client.Read(addrKey)
		if err != nil {
			log.WithError(err).Errorf("could not connect and read %s", addrKey)
			return err
		}

		if !ok {
			log.Errorf("account not found %s", addrKey)
			continue
		}

		namespace := strings.Split(addrKey, "0x")[0]
		if namespace == "" {
			namespace = multitenancy.DefaultTenant
		}

		acc, err := vault.ETHImportAccount(namespace, privKey)
		if err != nil {
			log.WithError(err).Errorf("Could not connect read %s", addrKey)
			return err
		}

		log.WithField("address", acc.Address).WithField("namespace", acc.Namespace).
			Infof("Account was imported successfully")
	}

	log.Info("Accounts have been imported successfully")

	return nil
}
