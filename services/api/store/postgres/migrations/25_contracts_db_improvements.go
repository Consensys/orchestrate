package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeRemoveFkeysContracts(db migrations.DB) error {
	log.Debug("Applying improvements contract tables and add index on idempotency-key...")
	_, err := db.Exec(`
CREATE UNIQUE INDEX tags_repository_name_idx ON tags (repository_id, (lower(name)));
ALTER TABLE repositories drop constraint repositories_name_key;
CREATE UNIQUE INDEX repositories_name_idx ON repositories ((lower(name)));

ALTER TABLE transactions
	ADD COLUMN contract_name TEXT,
	ADD COLUMN contract_tag TEXT;

CREATE INDEX transactions_hash_idem_idx on transaction_requests (idempotency_key);
`)
	if err != nil {
		return err
	}
	log.Info("Applied improvements on contract tables and added index on idempotency-key successfully")

	return nil
}

func downgradeRemoveFkeysContracts(db migrations.DB) error {
	log.Debug("Undoing improvements on contract tables and removing index on idempotency-key...")
	_, err := db.Exec(`
ALTER TABLE transactions
	DROP COLUMN contract_name,
	DROP COLUMN contract_tag;

DROP INDEX repositories_name_idx;
ALTER TABLE repositories ADD CONSTRAINT repositories_name_key UNIQUE (name);
DROP INDEX tags_repository_name_idx;

DROP INDEX transactions_hash_idem_idx;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded improvements on contract tables and removed index on idempotency-key")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeRemoveFkeysContracts, downgradeRemoveFkeysContracts)
}
