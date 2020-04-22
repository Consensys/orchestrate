package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createPrivateTxManagers(db migrations.DB) error {
	log.Debug("Creating tables...")
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
		log.WithError(err).Error("Could not create tables")
		return err
	}
	log.Info("Created tables")

	return nil
}

func dropPrivateTxManagers(db migrations.DB) error {
	log.Debug("Dropping tables")
	_, err := db.Exec(`
DROP TABLE private_tx_managers;

DROP TYPE priv_chain_type;
`)
	if err != nil {
		log.WithError(err).Error("Could not drop tables")
		return err
	}
	log.Info("Dropped tables")

	return nil
}

func init() {
	Collection.MustRegisterTx(createPrivateTxManagers, dropPrivateTxManagers)
}
