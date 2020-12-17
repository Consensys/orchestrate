package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createAccountsTable(db migrations.DB) error {
	log.Debug("Creating accounts table...")
	_, err := db.Exec(`
CREATE TABLE accounts (
	id SERIAL PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    alias TEXT,
    address CHAR(42) NOT NULL,
    public_key TEXT NOT NULL,
    compressed_public_key TEXT,
    active BOOLEAN default true,
	attributes JSONB,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL
);

CREATE UNIQUE INDEX account_unique_alias_idx ON accounts (tenant_id, alias) WHERE alias IS NOT NULL;
CREATE UNIQUE INDEX account_unique_address_idx ON accounts (address);

CREATE TRIGGER accounts_trigger
	BEFORE UPDATE ON accounts
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();
`)
	if err != nil {
		log.WithError(err).Error("Could not create accounts table")
		return err
	}
	log.Info("Created accounts table")

	return nil
}

func dropAccountsTable(db migrations.DB) error {
	log.Debug("Dropping accounts table")
	_, err := db.Exec(`
DROP TRIGGER accounts_trigger ON accounts;

DROP TABLE accounts;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop account table")
		return err
	}
	log.Info("Dropped accounts table")

	return nil
}

func init() {
	Collection.MustRegisterTx(createAccountsTable, dropAccountsTable)
}
