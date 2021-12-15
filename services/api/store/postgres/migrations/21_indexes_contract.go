package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeAddContractIndexes(db migrations.DB) error {
	log.Debug("Applying contract indexes...")
	_, err := db.Exec(`
CREATE INDEX artifacts_codehash_idx on artifacts (codehash);

CREATE INDEX codehashes_codehash_address_idx on codehashes (codehash, address);
`)

	if err != nil {
		return err
	}
	log.Info("Applied contract indexes")

	return nil
}

func downgradeAddContractIndexes(db migrations.DB) error {
	log.Debug("Downgrading contract indexes...")
	_, err := db.Exec(`
DROP INDEX artifacts_codehash_idx;

DROP INDEX codehashes_codehash_address_idx;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded contract indexes")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeAddContractIndexes, downgradeAddContractIndexes)
}
