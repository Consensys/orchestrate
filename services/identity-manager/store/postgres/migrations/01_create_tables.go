package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createContextTable(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
CREATE TABLE identities (
	id SERIAL PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    alias TEXT NOT NULL,
    address CHAR(42) NOT NULL,
    public_key TEXT NOT NULL,
    compressed_public_key TEXT,
    active BOOLEAN default true,
	attributes JSONB,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
	UNIQUE(tenant_id, alias)
);

CREATE OR REPLACE FUNCTION updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.updated_at = (now() at time zone 'utc');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER identities_trigger
	BEFORE UPDATE ON identities
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();
`)
	if err != nil {
		log.WithError(err).Error("Could not create tables")
		return err
	}
	log.Info("Created tables")

	return nil
}

func dropContextTable(db migrations.DB) error {
	log.Debug("Dropping tables")
	_, err := db.Exec(`
DROP TRIGGER identities_trigger ON identities;

DROP FUNCTION updated();

DROP TABLE identities;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop tables")
		return err
	}
	log.Info("Dropped tables")

	return nil
}

func init() {
	Collection.MustRegisterTx(createContextTable, dropContextTable)
}
