package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func createPrivateTxManagers(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
CREATE TABLE private_tx_managers (
	id UUID PRIMARY KEY,
	chain_uuid VARCHAR(66) NOT NULL,
	url TEXT NOT NULL,
	type VARCHAR(66) NOT NULL,
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
