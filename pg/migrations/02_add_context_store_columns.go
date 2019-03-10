package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func addColumnsOnContext(db migrations.DB) error {
	log.Debugf("Adding columns on table %q...", "context")

	// Remark: you will note that we consider that chain ID should be max a uint256
	_, err := db.Exec(
		`CREATE TYPE status AS ENUM ('pending', 'mined');
ALTER TABLE context 
	ADD COLUMN id SERIAL PRIMARY KEY, 
	ADD COLUMN chain_id varchar(66), 
	ADD COLUMN tx_hash char(66), 
	ADD COLUMN trace_id uuid, 
	ADD COLUMN status status, 
	ADD COLUMN sending_time timestamptz default (now() at time zone 'utc'), 
	ADD COLUMN mining_time timestamptz, 
	ADD COLUMN trace bytea;
CREATE UNIQUE INDEX tx_idx ON context (chain_id, tx_hash);
CREATE UNIQUE INDEX trace_idx ON context (trace_id);`,
	)

	if err != nil {
		return err
	}

	log.Infof("Added columns %q on table %q", []string{"id", "chain_id", "tx_hash", "trace_id", "status", "sending_time", "mining_time", "trace"}, "context")

	return nil
}

func dropColumnsOnContext(db migrations.DB) error {
	log.Debugf("Removing columns on table %q...", "context")

	_, err := db.Exec(
		`DROP INDEX tx_idx, trace_idx;
ALTER TABLE context 
	DROP COLUMN id, 
	DROP COLUMN chain_id, 
	DROP COLUMN tx_hash, 
	DROP COLUMN trace_id, 
	DROP COLUMN status, 
	DROP COLUMN sending_time, 
	DROP COLUMN mining_time, 
	DROP COLUMN trace;
DROP TYPE status`,
	)

	if err != nil {
		return err
	}

	log.Infof("Removed columns %q on table %q", []string{"id", "chain_id", "tx_hash", "trace_id", "status", "sending_time", "mining_time", "trace"}, "context")

	return nil
}

func init() {
	Collections.MustRegisterTx(addColumnsOnContext, dropColumnsOnContext)
}
