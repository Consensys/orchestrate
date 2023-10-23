package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeAddChainHeaders(db migrations.DB) error {
	log.Debug("Applying adding chain headers...")
	_, err := db.Exec(`
ALTER TABLE chains
	ADD COLUMN headers JSONB;
`)
	if err != nil {
		return err
	}
	log.Info("Applied adding chain headers")

	return nil
}

func downgradeAddChainHeaders(db migrations.DB) error {
	log.Debug("Downgrading adding chain headers...")
	_, err := db.Exec(`
ALTER TABLE chains
	DROP COLUMN headers;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded adding chain headers")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeAddChainHeaders, downgradeAddChainHeaders)
}
