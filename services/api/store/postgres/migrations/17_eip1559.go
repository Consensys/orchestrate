package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeEIP1559Support(db migrations.DB) error {
	log.Debug("Applying eip-1559 support...")
	_, err := db.Exec(`
ALTER TABLE transactions
	ADD COLUMN gas_fee_cap BIGINT,
	ADD COLUMN gas_tip_cap BIGINT,
	ADD COLUMN tx_type TEXT,
	ADD COLUMN access_list JSONB;
`)
	if err != nil {
		return err
	}
	log.Info("Apply adding chain labels")

	return nil
}

func downgradeEIP1559Support(db migrations.DB) error {
	log.Debug("Downgrading eip-1559 support...")
	_, err := db.Exec(`
ALTER TABLE transactions
	DROP COLUMN gas_fee_cap,
	DROP COLUMN gas_tip_cap,
	DROP COLUMN tx_type,
	DROP COLUMN access_list;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded eip-1559 support")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeEIP1559Support, downgradeEIP1559Support)
}
