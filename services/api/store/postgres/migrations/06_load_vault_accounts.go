package migrations

import (
	"context"
	"fmt"
	"strings"

	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func loadVaultAccounts(db migrations.DB) error {
	log.Debug("Loading accounts from Vault...")
	ctx := context.Background()

	storeName := qkm.GlobalStoreName()
	client := qkm.GlobalClient()
	if client == nil {
		log.Warnf("loading vault accounts ignored. Missing key-manager client")
		return nil
	}

	accounts, err := client.ListEthAccounts(ctx, storeName, 0, 0)
	if err != nil {
		log.WithError(err).Errorf("could not get list of accounts")
		return err
	}

	var queryInsertItems []string
	for _, accountID := range accounts {
		acc, err2 := client.GetEthAccount(ctx, storeName, accountID)
		if err2 != nil {
			log.WithField("account_id", accountID).WithError(err2).Error("Could not get account")
			return err2
		}

		tenantIDs := strings.Split(acc.Tags[qkm.TagIDAllowedTenants], qkm.TagSeparatorAllowedTenants)
		for _, tenantID := range tenantIDs {
			queryInsertItems = append(queryInsertItems, fmt.Sprintf("('%s', '%s', '%s', '%s', '{\"source\": \"kv-v2\"}')",
				tenantID,
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
