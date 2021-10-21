package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeAddChainLabels(db migrations.DB) error {
	log.Debug("Applying adding chain labels...")
	_, err := db.Exec(`
ALTER TABLE chains
	ADD COLUMN labels JSONB;
`)
	if err != nil {
		return err
	}
	log.Info("Applied adding chain labels")

	return nil
}

func downgradeAddChainLabels(db migrations.DB) error {
	log.Debug("Downgrading adding chain labels...")
	_, err := db.Exec(`
ALTER TABLE chains
	DROP COLUMN labels;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded adding chain labels")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeAddChainLabels, downgradeAddChainLabels)
}
