package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeDeprecateOrion(db migrations.DB) error {
	log.Debug("Applying deprecate Orion...")
	_, err := db.Exec(`
ALTER TYPE job_type RENAME VALUE 'eth://orion/markingTransaction' TO 'eth://eea/markingTransaction';
ALTER TYPE job_type RENAME VALUE 'eth://orion/eeaTransaction' TO 'eth://eea/privateTransaction';

ALTER TYPE priv_chain_type RENAME VALUE 'Orion' TO 'EEA';
`)

	if err != nil {
		return err
	}
	log.Info("Applied deprecated Orion")

	return nil
}

func downgradeDeprecateOrion(db migrations.DB) error {
	log.Debug("Downgrading deprecate Orion...")
	_, err := db.Exec(`
ALTER TYPE job_type RENAME VALUE 'eth://eea/markingTransaction' TO 'eth://orion/markingTransaction';
ALTER TYPE job_type RENAME VALUE 'eth://eea/privateTransaction' TO 'eth://orion/eeaTransaction';

ALTER TYPE priv_chain_type RENAME VALUE 'EEA' TO 'Orion';
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded deprecate Orion")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeDeprecateOrion, downgradeDeprecateOrion)
}
