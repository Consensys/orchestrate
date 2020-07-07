package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func addChainIDColumn(db migrations.DB) error {
	log.Debugf("Adding chainID column on table %q...", "chains")

	_, err := db.Exec(`
ALTER TABLE chains
	ADD COLUMN chain_id BIGINT NOT NULL;
	`)

	if err != nil {
		return err
	}

	log.Infof("Added chainID columns on table %q", "chains")

	return nil
}

func dropChainIDColumn(db migrations.DB) error {
	log.Debugf("Removing chainID chainID on table %q...", "chains")

	_, err := db.Exec(`
ALTER TABLE chains 
	DROP COLUMN chain_id;
	`)

	if err != nil {
		return err
	}

	log.Infof("Removed chainID column on table %q", "chains")

	return nil
}

func init() { Collection.MustRegisterTx(addChainIDColumn, dropChainIDColumn) }
