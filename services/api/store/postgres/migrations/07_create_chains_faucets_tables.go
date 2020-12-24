package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createChainRegistryTables(db migrations.DB) error {
	log.Debug("Creating chains and faucets tables...")

	_, err := db.Exec(`
CREATE TABLE chains (
	uuid UUID PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	tenant_id VARCHAR(66) NOT NULL,
	urls TEXT[] NOT NULL,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	listener_depth INTEGER,
	listener_current_block BIGINT,
	listener_starting_block BIGINT,
	listener_back_off_duration VARCHAR(66) NOT NULL,
	listener_external_tx_enabled BOOLEAN DEFAULT false NOT NULL
);
CREATE UNIQUE INDEX ON chains (tenant_id, name);

CREATE TABLE faucets (
	uuid UUID PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	tenant_id VARCHAR(66) NOT NULL,
	
	chain_rule VARCHAR(255), 
	creditor_account CHAR(42) NOT NULL,
	max_balance BIGINT,
	amount BIGINT,
	cooldown VARCHAR(66),
	
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL, 
	updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL
);

CREATE TRIGGER chain_trigger
	BEFORE UPDATE ON chains
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER faucet_trigger
	BEFORE UPDATE ON faucets
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();
`)
	if err != nil {
		log.WithError(err).Error("Could not create chains and faucets tables")
		return err
	}
	log.Info("Created chains and faucets tables")

	return nil
}

func dropChainRegistryTables(db migrations.DB) error {
	log.Debug("Dropping chains and faucets tables")

	_, err := db.Exec(`
DROP TRIGGER chain_trigger ON chains;
DROP TABLE chains;

DROP TRIGGER faucet_trigger ON faucets;
DROP TABLE faucets;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop faucets and chains tables")
		return err
	}
	log.Info("Dropped chains and faucets tables")

	return nil
}

func init() {
	Collection.MustRegisterTx(createChainRegistryTables, dropChainRegistryTables)
}
