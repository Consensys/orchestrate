package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createContextTable(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
CREATE TABLE chains (
	uuid UUID PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	tenant_id VARCHAR(66) NOT NULL,
	urls TEXT[] NOT NULL,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	listener_depth INTEGER,
	listener_block_position BIGINT,
	listener_from_block BIGINT,
	listener_back_off_duration VARCHAR(66) NOT NULL
);
CREATE UNIQUE INDEX ON chains (tenant_id, name);

CREATE OR REPLACE FUNCTION chain_updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.updated_at = (now() at time zone 'utc');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER chain_trigger
	BEFORE UPDATE ON chains
	FOR EACH ROW 
	EXECUTE PROCEDURE chain_updated();
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
DROP TRIGGER chain_trigger ON chains;
DROP FUNCTION chain_updated();
DROP TABLE chains;
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
