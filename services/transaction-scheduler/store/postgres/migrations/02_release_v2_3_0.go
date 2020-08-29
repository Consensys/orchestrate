package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upMigration02(db migrations.DB) error {
	log.Debug("Applying migration 02...")
	_, err := db.Exec(`
ALTER TABLE transactions 
	ALTER COLUMN value TYPE TEXT using value::text,
	ALTER COLUMN nonce TYPE TEXT using nonce::text,
	ALTER COLUMN gas_price TYPE TEXT using gas_price::text;
`)
	if err != nil {
		log.WithError(err).Error("Could not apply migration 02")
		return err
	}
	log.Info("Migration completed")

	return nil
}

func downMigration02(db migrations.DB) error {
	log.Debug("Rollback migration 02")
	_, err := db.Exec(`
ALTER TABLE transactions 
	ALTER COLUMN value TYPE BIGINT USING value::bigint,
	ALTER COLUMN nonce TYPE INTEGER USING nonce::int,
	ALTER COLUMN gas_price TYPE BIGINT USING gas_price::bigint;
`)
	if err != nil {
		log.WithError(err).Error("Could not apply rollback")
		return err
	}

	log.Info("Rollback completed")
	return nil
}

func init() {
	Collection.MustRegisterTx(upMigration02, downMigration02)
}
