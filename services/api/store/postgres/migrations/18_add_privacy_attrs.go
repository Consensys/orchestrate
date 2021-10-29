package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeNewPrivacyAttrSupport(db migrations.DB) error {
	log.Debug("Applying new privacy attributes...")
	_, err := db.Exec(`
ALTER TABLE transactions
	ADD COLUMN mandatory_for TEXT [],
	ADD COLUMN privacy_flag int;
`)

	if err != nil {
		return err
	}
	log.Info("Applied new privacy attributes")

	return nil
}

func downgradePrivacyFlagSupport(db migrations.DB) error {
	log.Debug("Downgrading new privacy attributes...")
	_, err := db.Exec(`
ALTER TABLE transactions
	DROP COLUMN mandatory_for,
	DROP COLUMN privacy_flag;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded new privacy attributes")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeNewPrivacyAttrSupport, downgradePrivacyFlagSupport)
}
