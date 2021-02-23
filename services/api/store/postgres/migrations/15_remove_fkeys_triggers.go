package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeRemoveFkeysTrigger(db migrations.DB) error {
	log.Debug("Applying removing fkeys and triggers...")
	_, err := db.Exec(`
DROP FUNCTION updated CASCADE;

DROP FUNCTION job_log_updated CASCADE;

ALTER TABLE jobs
	DROP CONSTRAINT jobs_schedule_id_fkey;

ALTER TABLE jobs
	DROP CONSTRAINT jobs_transaction_id_fkey;

ALTER TABLE jobs
	DROP CONSTRAINT jobs_chain_uuid_fkey;

ALTER TABLE transaction_requests
	DROP CONSTRAINT transaction_requests_schedule_id_fkey;

ALTER TABLE logs
	DROP CONSTRAINT logs_job_id_fkey;

CREATE INDEX jobs_schedule_id_idx on jobs (schedule_id);
`)
	if err != nil {
		return err
	}
	log.Info("Applying removing fkeys and triggers successfully")

	return nil
}

func downgradeRemoveFkeysTrigger(db migrations.DB) error {
	log.Debug("Undoing removing fkeys and triggers...")
	_, err := db.Exec(`
CREATE OR REPLACE FUNCTION updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.updated_at = (now() at time zone 'utc');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE or REPLACE FUNCTION job_log_updated() RETURNS trigger AS
	$$
	BEGIN
	  UPDATE jobs SET updated_at = (now() at time zone 'utc') WHERE id = NEW.job_id;
	  RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER job_trigger
	BEFORE UPDATE ON jobs
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER update_parent_job
	AFTER INSERT OR UPDATE ON logs
	FOR EACH ROW
	EXECUTE PROCEDURE job_log_updated();

CREATE TRIGGER accounts_updated_trigger
	BEFORE UPDATE ON accounts
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER chain_trigger
	BEFORE UPDATE ON chains
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER faucet_trigger
	BEFORE UPDATE ON faucets
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER transactions_updated_trigger
	BEFORE UPDATE ON transactions
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER accounts_trigger
	BEFORE UPDATE ON accounts
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

ALTER TABLE jobs
	ADD CONSTRAINT jobs_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES schedules (id) ON DELETE CASCADE;

ALTER TABLE jobs
	ADD CONSTRAINT jobs_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions (id) ON DELETE RESTRICT;

ALTER TABLE jobs
	ADD CONSTRAINT jobs_chain_uuid_fkey FOREIGN KEY (chain_uuid) REFERENCES chains (uuid) ON DELETE SET NULL;

ALTER TABLE transaction_requests
	ADD CONSTRAINT transaction_requests_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES schedules (id) ON DELETE SET NULL;

ALTER TABLE logs
	ADD CONSTRAINT logs_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE;

DROP INDEX jobs_schedule_id_idx;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded removing fkeys and triggers")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeRemoveFkeysTrigger, downgradeRemoveFkeysTrigger)
}
