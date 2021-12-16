package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeAccountStoreID(db migrations.DB) error {
	log.Debug("Applying adding storeID to accounts...")
	_, err := db.Exec(`
ALTER TABLE accounts
	ADD COLUMN store_id TEXT;
`)
	if err != nil {
		return err
	}
	log.Info("Applied adding storeID to accounts")

	return nil
}

func downgradeAccountStoreID(db migrations.DB) error {
	log.Debug("Downgrading adding storeID to accounts...")
	_, err := db.Exec(`
ALTER TABLE accounts
	DROP COLUMN store_id;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded adding storeID to accounts")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeAccountStoreID, downgradeAccountStoreID)
}
