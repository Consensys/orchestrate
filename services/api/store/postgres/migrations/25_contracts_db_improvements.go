package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeRemoveFkeysContracts(db migrations.DB) error {
	log.Debug("Applying improvements contract tables...")
	_, err := db.Exec(`
CREATE UNIQUE INDEX tags_repository_name_idx ON tags (repository_id, (lower(name)));
ALTER TABLE repositories drop constraint repositories_name_key;
CREATE UNIQUE INDEX repositories_name_idx ON repositories ((lower(name)));

ALTER TABLE transactions
	ADD COLUMN contract_name TEXT,
	ADD COLUMN contract_tag TEXT;
`)
	if err != nil {
		return err
	}
	log.Info("Applied improvements on contract tables successfully")

	return nil
}

func downgradeRemoveFkeysContracts(db migrations.DB) error {
	log.Debug("Undoing improvements on contract tables...")
	_, err := db.Exec(`
ALTER TABLE transactions
	DROP COLUMN contract_name,
	DROP COLUMN contract_tag;

DROP INDEX repositories_name_idx;
ALTER TABLE repositories ADD CONSTRAINT repositories_name_key UNIQUE (name);
DROP INDEX tags_repository_name_idx;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded improvements on contract tables")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeRemoveFkeysContracts, downgradeRemoveFkeysContracts)
}
