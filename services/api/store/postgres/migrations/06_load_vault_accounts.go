package migrations

import (
	"context"
	"fmt"
	"strings"

	keymanagerclient "github.com/ConsenSys/orchestrate/services/key-manager/client"
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func loadVaultAccounts(db migrations.DB) error {
	log.Debug("Loading accounts from Vault...")
	ctx := context.Background()

	client := keymanagerclient.GlobalClient()
	if client == nil {
		log.Warnf("loading vault accounts ignored. Missing key-manager client")
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
			log.WithField("namespace", namespace).WithError(err2).Errorf("Could not get list of accounts")
			return err2
		}

		for _, addr := range accounts {
			acc, err2 := client.ETHGetAccount(ctx, addr, namespace)
			if err2 != nil {
				log.WithField("namespace", namespace).WithField("address", addr).
					WithError(err2).Error("Could not get account")
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
			log.WithError(err).Error("Could not import accounts")
			return err
		}
	}

	log.WithField("accounts", len(queryInsertItems)).Info("accounts imported successfully")
	return nil
}

func dropVaultAccounts(db migrations.DB) error {
	log.Debug("Dropping tables")
	_, err := db.Exec(`
DELETE FROM accounts WHERE attributes @> '{"source": "kv-v2"}'
`)
	if err != nil {
		log.WithError(err).Error("Could not drop vault accounts")
		return err
	}
	log.Info("Dropped vault accounts")

	return nil
}

func init() {
	Collection.MustRegisterTx(loadVaultAccounts, dropVaultAccounts)
}
