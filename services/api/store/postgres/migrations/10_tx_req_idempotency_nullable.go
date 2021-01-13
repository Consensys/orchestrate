package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upMigration10(db migrations.DB) error {
	log.Debug("Marking idempotencyKey nullable ...")
	_, err := db.Exec(`
ALTER TABLE transaction_requests 
	ALTER COLUMN idempotency_key DROP NOT NULL;
`)
	if err != nil {
		log.WithError(err).Error("Could not apply migration")
		return err
	}
	log.Info("Migration 10 completed")

	return nil
}

func downMigration10(db migrations.DB) error {
	log.Debug("Rollback marking idempotencyKey nullable")
	_, err := db.Exec(`
ALTER TABLE transaction_requests 
	ALTER COLUMN idempotency_key SET NOT NULL;
`)
	if err != nil {
		log.WithError(err).Error("Could not apply rollback")
		return err
	}

	log.Info("Rollback migration 10 completed")
	return nil
}

func init() {
	Collection.MustRegisterTx(upMigration10, downMigration10)
}
