package migrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
	keyManagerClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

func loadVaultIdentities(db migrations.DB) error {
	log.Debug("Loading identities from Vault...")
	ctx := context.Background()

	client := keyManagerClient.GlobalClient()
	if client == nil {
		log.Warnf("migration is ignored. Missing key-manager client")
		return nil
	}

	namespaces, err := client.ETHListNamespaces(ctx)
	if err != nil {
		log.WithError(err).Errorf("could not get list of namespaces")
		return err
	}

	var queryInsertItems []string
	for _, namespace := range namespaces {
		accounts, err2 := client.ETHListAccounts(ctx, namespace)
		if err2 != nil {
			log.WithField("namespace", namespace).WithError(err2).Errorf("Could not get list of account")
			return err2
		}

		for _, addr := range accounts {
			acc, err2 := client.ETHGetAccount(ctx, addr, namespace)
			if err2 != nil {
				log.WithField("namespace", namespace).WithField("address", addr).
					WithError(err2).Errorf("Could not get account")
				return err2
			}

			queryInsertItems = append(queryInsertItems, fmt.Sprintf("('%s', '%s', '%s', '%s', '{\"source\": \"kv-v2\"}')",
				acc.Namespace,
				acc.Address,
				acc.PublicKey,
				acc.CompressedPublicKey,
			))
		}
	}

	if len(queryInsertItems) > 0 {
		_, err = db.Exec("INSERT INTO accounts (tenant_id, address, public_key, compressed_public_key, attributes) VALUES " +
			strings.Join(queryInsertItems, ", ") + " on conflict do nothing")
		if err != nil {
			log.WithError(err).Error("Could not import identities")
			return err
		}
	}

	log.Infof("%d Identities imported successfully", len(queryInsertItems))

	return nil
}

func dropVaultIdentities(db migrations.DB) error {
	log.Debug("Dropping tables")
	_, err := db.Exec(`
DELETE FROM accounts WHERE attributes @> '{"source": "kv-v2"}'
`)
	if err != nil {
		log.WithError(err).Error("Could not drop tables")
		return err
	}
	log.Info("Dropped tables")

	return nil
}

func init() {
	Collection.MustRegisterTx(loadVaultIdentities, dropVaultIdentities)
}
