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
	listener_depth INTEGER,
	listener_block_position BIGINT,
	listener_from_block BIGINT,
	listener_back_off_duration VARCHAR(66) NOT NULL,
	listener_external_tx_enabled BOOLEAN DEFAULT false NOT NULL,
	
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL
);
CREATE UNIQUE INDEX ON chains (tenant_id, name);

CREATE TABLE faucets (
	uuid UUID PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	tenant_id VARCHAR(66) NOT NULL,
	
	chain_uuid UUID REFERENCES chains,
	max_balance BIGINT,
	creditor_account_address CHAR(42) NOT NULL,
	cooldown VARCHAR(66),
	
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL
);
CREATE UNIQUE INDEX ON faucets (tenant_id, name);

CREATE OR REPLACE FUNCTION updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.updated_at = (now() at time zone 'utc');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER chain_trigger
	BEFORE UPDATE ON chains
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
DROP TRIGGER chain_trigger ON chains;
DROP FUNCTION updated();
DROP TABLE faucets;
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
