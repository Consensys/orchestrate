package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upMigration04(db migrations.DB) error {
	log.Debug("Applying migration 04...")
	_, err := db.Exec(`
UPDATE schedules
	SET uuid=transaction_requests.uuid
	FROM transaction_requests
	WHERE schedules.id = transaction_requests.id;

ALTER TABLE transaction_requests
	DROP COLUMN uuid;
`)
	if err != nil {
		log.WithError(err).Error("Could not apply migration 04")
		return err
	}
	log.Info("Migration 04 completed")

	return nil
}

func downMigration04(db migrations.DB) error {
	log.Debug("Rollback migration 04")
	_, err := db.Exec(`
ALTER TABLE transaction_requests
	ADD COLUMN uuid UUID;

UPDATE transaction_requests
	SET uuid=schedules.uuid
	FROM schedules
	WHERE schedules.id = transaction_requests.id;

ALTER TABLE transaction_requests
	ALTER COLUMN uuid SET NOT NULL;
`)
	if err != nil {
		log.WithError(err).Error("Could not apply rollback")
		return err
	}

	log.Info("Rollback migration 04 completed")
	return nil
}

func init() {
	Collection.MustRegisterTx(upMigration04, downMigration04)
}
