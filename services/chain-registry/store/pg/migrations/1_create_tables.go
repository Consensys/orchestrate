package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func createContextTable(db migrations.DB) error {
	log.Debug("Creating tables...")
	_, err := db.Exec(`
CREATE TABLE config (
	id SERIAL PRIMARY KEY,
	name VARCHAR(66) NOT NULL,
	tenant_id VARCHAR(66) DEFAULT 'default',
	config_type INTEGER NOT NULL,
	config JSONB NOT NULL
);
CREATE INDEX ON config (tenant_id, name, config_type);
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
DROP TABLE config;
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
