package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upMigration03(db migrations.DB) error {
	log.Debug("Applying migration 03...")
	_, err := db.Exec(`
ALTER TABLE jobs
	ADD COLUMN updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL;

ALTER TABLE jobs
	RENAME COLUMN annotations TO internal_data;

CREATE OR REPLACE FUNCTION updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.updated_at = (now() at time zone 'utc');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER job_trigger
	BEFORE UPDATE ON jobs
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();
`)
	if err != nil {
		log.WithError(err).Error("Could not apply migration 03")
		return err
	}
	log.Info("Migration 03 completed")

	return nil
}

func downMigration03(db migrations.DB) error {
	log.Debug("Rollback migration 03")
	_, err := db.Exec(`
ALTER TABLE jobs
	DROP COLUMN updated_at;

ALTER TABLE jobs
  RENAME COLUMN internal_data TO annotations;

DROP TRIGGER job_trigger ON jobs;
DROP FUNCTION updated();
`)
	if err != nil {
		log.WithError(err).Error("Could not apply rollback")
		return err
	}

	log.Info("Rollback migration 03 completed")
	return nil
}

func init() {
	Collection.MustRegisterTx(upMigration03, downMigration03)
}
