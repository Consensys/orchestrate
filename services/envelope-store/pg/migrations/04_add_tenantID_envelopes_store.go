package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func addTenantIDColumnsOnEnvelopeStore(db migrations.DB) error {
	log.Debugf("Adding tenantID columns on table %q...", "envelopes")

	// Remark: you will note that we consider that tenant_id should be max 255 characters
	_, err := db.Exec(`
ALTER TABLE envelopes
	ADD COLUMN tenant_id varchar(255) NOT NULL DEFAULT '_';
	`)

	if err != nil {
		return err
	}

	log.Infof("Added tenantID columns on table %q", "envelopes")

	return nil
}

func dropTenantIDColumnsOnEnvelopeStore(db migrations.DB) error {
	log.Debugf("Removing columns on table %q...", "envelopes")

	_, err := db.Exec(`
ALTER TABLE envelopes 
	DROP COLUMN tenant_id;
	`)

	if err != nil {
		return err
	}

	log.Infof("Removed tenantID columns on table %q", "envelopes")

	return nil
}

func init() {
	Collection.MustRegisterTx(addTenantIDColumnsOnEnvelopeStore, dropTenantIDColumnsOnEnvelopeStore)
}
