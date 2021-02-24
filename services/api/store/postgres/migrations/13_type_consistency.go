package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeTypesConsistency(db migrations.DB) error {
	log.Debug("Applying database type consistency...")
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

ALTER TABLE chains
	ALTER COLUMN name TYPE TEXT,
	ALTER COLUMN tenant_id TYPE TEXT;

CREATE TYPE job_type AS ENUM ('eth://ethereum/transaction', 'eth://ethereum/rawTransaction', 'eth://orion/markingTransaction', 'eth://orion/eeaTransaction', 'eth://tessera/markingTransaction', 'eth://tessera/privateTransaction');

ALTER TABLE jobs
	ALTER COLUMN chain_uuid DROP NOT NULL;

ALTER TABLE jobs
	ALTER COLUMN type TYPE job_type using type::job_type;

CREATE TYPE job_status AS ENUM ('CREATED', 'STARTED', 'PENDING', 'MINED', 'NEVER_MINED', 'RESENDING', 'STORED', 'RECOVERING', 'WARNING', 'FAILED');

ALTER TABLE logs
	ALTER COLUMN status TYPE job_status using status::job_status;
`)
	if err != nil {
		return err
	}
	log.Info("Database union refactor completed")

	return nil
}

func downgradeTypesConsistency(db migrations.DB) error {
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

ALTER TABLE chains
	ALTER COLUMN name TYPE VARCHAR(66),
	ALTER COLUMN tenant_id TYPE VARCHAR(66);

ALTER TABLE jobs
	ALTER COLUMN type TYPE TEXT;

DROP TYPE job_type;

ALTER TABLE logs
	ALTER COLUMN status TYPE TEXT;

DROP TYPE job_status;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded database refactor")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeTypesConsistency, downgradeTypesConsistency)
}
