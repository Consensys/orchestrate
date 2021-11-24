package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeOwnerSupport(db migrations.DB) error {
	log.Debug("Applying ownership support...")
	_, err := db.Exec(`
ALTER TABLE chains 
	ADD COLUMN owner_id TEXT; 

ALTER TABLE schedules 
	ADD COLUMN owner_id TEXT;

ALTER TABLE accounts 
	ADD COLUMN owner_id TEXT;

DROP INDEX account_unique_alias_idx;
CREATE UNIQUE INDEX account_unique_alias_idx ON accounts (tenant_id, owner_id, alias) WHERE alias IS NOT NULL;

DROP INDEX chains_tenant_id_name_idx;
CREATE UNIQUE INDEX chains_unique_name_idx ON chains (tenant_id, owner_id, name);
`)

	if err != nil {
		return err
	}
	log.Info("Applied ownership support")

	return nil
}

func downgradeOwnerSupport(db migrations.DB) error {
	log.Debug("Downgrading ownership support...")
	_, err := db.Exec(`
DROP INDEX account_unique_alias_idx;
CREATE UNIQUE INDEX account_unique_alias_idx ON accounts (tenant_id, alias) WHERE alias IS NOT NULL;

DROP INDEX chains_unique_name_idx;
CREATE UNIQUE INDEX chains_tenant_id_name_idx ON chains (tenant_id, name);

ALTER TABLE chains 
	DROP COLUMN owner_id;

ALTER TABLE schedules 
	DROP COLUMN owner_id;

ALTER TABLE accounts 
	DROP COLUMN owner_id;
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded ownership support")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeOwnerSupport, downgradeOwnerSupport)
}
