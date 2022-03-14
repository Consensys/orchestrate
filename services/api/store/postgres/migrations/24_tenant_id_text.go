package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeTenantFaucetTypeConsistency(db migrations.DB) error {
	log.Debug("Applying tenant faucet type consistency...")
	_, err := db.Exec(`
ALTER TABLE faucets
	ALTER COLUMN name TYPE TEXT,
	ALTER COLUMN chain_rule TYPE TEXT,
	ALTER COLUMN creditor_account TYPE TEXT,
	ALTER COLUMN tenant_id TYPE TEXT;
`)
	if err != nil {
		return err
	}
	log.Info("Database union refactor completed")

	return nil
}

func downgradeTenantFaucetTypeConsistency(db migrations.DB) error {
	log.Debug("Undoing tenant faucet type consistency")
	_, err := db.Exec(`
ALTER TABLE faucets
	ALTER COLUMN name TYPE VARCHAR(66),
	ALTER COLUMN chain_rule TYPE VARCHAR(66),
	ALTER COLUMN creditor_account TYPE VARCHAR(66),
	ALTER COLUMN tenant_id TYPE VARCHAR(66);
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded database refactor")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeTenantFaucetTypeConsistency, downgradeTenantFaucetTypeConsistency)
}
