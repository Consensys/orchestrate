package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func createContextTable(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
CREATE TABLE nodes (
	id SERIAL PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	tenant_id VARCHAR(66) NOT NULL,
	urls TEXT[] NOT NULL,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	listener_depth INTEGER NOT NULL,
	listener_block_position BIGINT NOT NULL,
	listener_from_block BIGINT NOT NULL,
	listener_back_off_duration VARCHAR(66) NOT NULL
);
CREATE UNIQUE INDEX ON nodes (tenant_id, name);

CREATE OR REPLACE FUNCTION node_updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.updated_at = (now() at time zone 'utc');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER node_trigger
	BEFORE UPDATE ON nodes
	FOR EACH ROW 
	EXECUTE PROCEDURE node_updated();
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
DROP TRIGGER node_trigger ON nodes;
DROP FUNCTION node_updated();
DROP TABLE nodes;
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
