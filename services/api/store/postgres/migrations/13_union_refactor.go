package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeUnionRefactor(db migrations.DB) error {
	log.Debug("Applying database union refactor...")
	_, err := db.Exec(`
ALTER TABLE accounts
	ALTER COLUMN address TYPE TEXT;

ALTER TABLE artifacts
	ALTER COLUMN codehash TYPE TEXT;

ALTER TABLE codehashes
	ALTER COLUMN codehash TYPE TEXT,
	ALTER COLUMN address TYPE TEXT,
	ALTER COLUMN chain_id TYPE TEXT;

ALTER TABLE methods
	ALTER COLUMN codehash TYPE TEXT,
	ALTER COLUMN selector TYPE TEXT;

ALTER TABLE events
	ALTER COLUMN codehash TYPE TEXT,
	ALTER COLUMN sig_hash TYPE TEXT;

ALTER TABLE repositories
	ALTER COLUMN name TYPE TEXT;

ALTER TABLE tags
	ALTER COLUMN name TYPE TEXT;

ALTER TABLE tags
	DROP CONSTRAINT tags_repository_id_fkey,
	ADD CONSTRAINT tags_repository_id_fkey FOREIGN KEY (repository_id) REFERENCES repositories (id) ON DELETE CASCADE;

ALTER TABLE tags
	DROP CONSTRAINT tags_artifact_id_fkey,
	ADD CONSTRAINT tags_artifact_id_fkey FOREIGN KEY (artifact_id) REFERENCES artifacts (id) ON DELETE RESTRICT;

ALTER TABLE chains
	ALTER COLUMN name TYPE TEXT,
	ALTER COLUMN tenant_id TYPE TEXT;

ALTER TABLE transaction_requests
	DROP CONSTRAINT transaction_requests_schedule_id_fkey,
	ADD CONSTRAINT transaction_requests_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES schedules (id) ON DELETE SET NULL;

CREATE TYPE job_type AS ENUM ('eth://ethereum/transaction', 'eth://ethereum/rawTransaction', 'eth://orion/markingTransaction', 'eth://orion/eeaTransaction', 'eth://tessera/markingTransaction', 'eth://tessera/privateTransaction');

ALTER TABLE jobs
	ALTER COLUMN type TYPE job_type using type::job_type,
	ADD CONSTRAINT jobs_chain_uuid_fkey FOREIGN KEY (chain_uuid) REFERENCES chains (uuid) ON DELETE RESTRICT;

ALTER TABLE jobs
	DROP CONSTRAINT jobs_schedule_id_fkey,
	ADD CONSTRAINT jobs_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES schedules (id) ON DELETE CASCADE;

ALTER TABLE jobs
	DROP CONSTRAINT jobs_transaction_id_fkey,
	ADD CONSTRAINT jobs_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions (id) ON DELETE RESTRICT;

CREATE TYPE job_status AS ENUM ('CREATED', 'STARTED', 'PENDING', 'MINED', 'NEVER_MINED', 'RESENDING', 'STORED', 'RECOVERING', 'WARNING', 'FAILED');

ALTER TABLE logs
	ALTER COLUMN status TYPE job_status using status::job_status;

CREATE TRIGGER accounts_updated_trigger
	BEFORE UPDATE ON accounts
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER chains_updated_trigger
	BEFORE UPDATE ON chains
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER faucets_updated_trigger
	BEFORE UPDATE ON faucets
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE TRIGGER transactions_updated_trigger
	BEFORE UPDATE ON transactions
	FOR EACH ROW 
	EXECUTE PROCEDURE updated();

CREATE or REPLACE FUNCTION job_log_updated() RETURNS trigger AS
	$$
	BEGIN
	  UPDATE jobs SET updated_at = (now() at time zone 'utc') WHERE id = NEW.job_id;
	  RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER update_parent_job
	AFTER INSERT OR UPDATE ON logs
	FOR EACH ROW
	EXECUTE PROCEDURE job_log_updated();
`)
	if err != nil {
		return err
	}
	log.Info("Database union refactor completed")

	return nil
}

func downgradeUnionRefactor(db migrations.DB) error {
	log.Debug("Undoing union refactor")
	_, err := db.Exec(`
ALTER TABLE accounts
	ALTER COLUMN address TYPE CHAR(42);

ALTER TABLE artifacts
	ALTER COLUMN codehash TYPE CHAR(66);

ALTER TABLE codehashes
	ALTER COLUMN codehash TYPE CHAR(66),
	ALTER COLUMN address TYPE char(42),
	ALTER COLUMN chain_id TYPE VARCHAR(66);

ALTER TABLE methods
	ALTER COLUMN codehash TYPE CHAR(66),
	ALTER COLUMN selector TYPE CHAR(10);

ALTER TABLE events
	ALTER COLUMN codehash TYPE CHAR(66),
	ALTER COLUMN sig_hash TYPE CHAR(66);

ALTER TABLE repositories
	ALTER COLUMN name TYPE VARCHAR(66);

ALTER TABLE tags
	ALTER COLUMN name TYPE VARCHAR(66);

ALTER TABLE tags
	DROP CONSTRAINT tags_repository_id_fkey,
	ADD CONSTRAINT tags_repository_id_fkey FOREIGN KEY (repository_id) REFERENCES repositories (id);

ALTER TABLE tags
	DROP CONSTRAINT tags_artifact_id_fkey,
	ADD CONSTRAINT tags_artifact_id_fkey FOREIGN KEY (artifact_id) REFERENCES artifacts (id);

ALTER TABLE chains
	ALTER COLUMN name TYPE TEXT,
	ALTER COLUMN tenant_id TYPE TEXT;

ALTER TABLE chains
	ALTER COLUMN name TYPE VARCHAR(66),
	ALTER COLUMN tenant_id TYPE VARCHAR(66);

ALTER TABLE transaction_requests
	DROP CONSTRAINT transaction_requests_schedule_id_fkey,
	ADD CONSTRAINT transaction_requests_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES schedules (id);

ALTER TABLE jobs
	ALTER COLUMN type TYPE TEXT,
	DROP CONSTRAINT jobs_chain_uuid_fkey;

ALTER TABLE jobs
	DROP CONSTRAINT jobs_schedule_id_fkey,
	ADD CONSTRAINT jobs_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES schedules (id);

ALTER TABLE jobs
	DROP CONSTRAINT jobs_transaction_id_fkey,
	ADD CONSTRAINT jobs_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions (id);

DROP TYPE job_type;

ALTER TABLE logs
	ALTER COLUMN status TYPE TEXT;

DROP TYPE job_status;

DROP TRIGGER accounts_updated_trigger ON accounts;
DROP TRIGGER chains_updated_trigger ON chains;
DROP TRIGGER faucets_updated_trigger ON faucets;
DROP TRIGGER transactions_updated_trigger ON transactions;
DROP TRIGGER update_parent_job on logs;

DROP FUNCTION job_log_updated() CASCADE;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded database refactor")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeUnionRefactor, downgradeUnionRefactor)
}
