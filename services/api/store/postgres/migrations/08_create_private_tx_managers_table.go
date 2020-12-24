package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createPrivateTxManagersTable(db migrations.DB) error {
	log.Debug("Creating private_tx_managers tables...")
	_, err := db.Exec(`

CREATE TYPE priv_chain_type AS ENUM ('Tessera', 'Orion');

CREATE TABLE private_tx_managers (
	uuid UUID PRIMARY KEY,
	chain_uuid UUID NOT NULL REFERENCES chains(uuid) ON DELETE CASCADE,
	url TEXT NOT NULL,
	type priv_chain_type NOT NULL,
	created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL
);
`)
	if err != nil {
		log.WithError(err).Error("Could not create private_tx_managers table")
		return err
	}
	log.Info("Created private_tx_managers table")

	return nil
}

func dropPrivateTxManagersTable(db migrations.DB) error {
	log.Debug("Dropping private_tx_managers tables")

	_, err := db.Exec(`
DROP TABLE private_tx_managers;
DROP TYPE priv_chain_type;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop private_tx_managers table")
		return err
	}
	log.Info("Dropped private_tx_managers table")

	return nil
}

func init() {
	Collection.MustRegisterTx(createPrivateTxManagersTable, dropPrivateTxManagersTable)
}
