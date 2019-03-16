package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func addColumnsOnTraceStore(db migrations.DB) error {
	log.Debugf("Adding columns on table %q...", "traces")

	// Remark: you will note that we consider that chain ID should be max a uint256
	_, err := db.Exec(
		`CREATE TYPE status AS ENUM ('stored', 'error', 'pending', 'mined');
ALTER TABLE traces 
	ADD COLUMN id serial PRIMARY KEY, 
	ADD COLUMN chain_id varchar(66) NOT NULL, 
	ADD COLUMN tx_hash char(66) NOT NULL, 
	ADD CONSTRAINT uni_tx UNIQUE (chain_id, tx_hash),
	ADD COLUMN trace_id uuid NOT NULL UNIQUE, 
	ADD COLUMN status status default 'stored' NOT NULL, 
	ADD COLUMN stored_at timestamptz default (now() at time zone 'utc') NOT NULL, 
	ADD COLUMN error_at timestamptz, 
	ADD COLUMN sent_at timestamptz, 
	ADD COLUMN mined_at timestamptz, 
	ADD COLUMN trace bytea NOT NULL;

CREATE OR REPLACE FUNCTION status_updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		CASE NEW.status
			WHEN 'error' THEN
				NEW.error_at = (now() at time zone 'utc');
			WHEN 'pending' THEN
				NEW.sent_at = (now() at time zone 'utc');
			WHEN 'mined' THEN
				NEW.mined_at = (now() at time zone 'utc');
			ELSE
		END CASE;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER status_trig 
	BEFORE INSERT OR UPDATE OF status ON traces
	FOR EACH ROW EXECUTE PROCEDURE status_updated();`,
	)

	if err != nil {
		return err
	}

	log.Infof("Added columns on table %q", "traces")

	return nil
}

func dropColumnsOnTraceStore(db migrations.DB) error {
	log.Debugf("Removing columns on table %q...", "traces")

	_, err := db.Exec(
		`DROP TRIGGER status_trig ON traces;
DROP FUNCTION status_updated();

ALTER TABLE traces 
	DROP COLUMN id, 
	DROP COLUMN chain_id, 
	DROP COLUMN tx_hash, 
	DROP COLUMN trace_id, 
	DROP COLUMN status, 
	DROP COLUMN stored_at, 
	DROP COLUMN error_at, 
	DROP COLUMN sent_at, 
	DROP COLUMN mined_at, 
	DROP COLUMN trace;
DROP TYPE status;`,
	)

	if err != nil {
		return err
	}

	log.Infof("Removed columns on table %q", "traces")

	return nil
}

func init() {
	Collections.MustRegisterTx(addColumnsOnTraceStore, dropColumnsOnTraceStore)
}
