package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func loadVaultAccounts(_ migrations.DB) error {
	log.Debug("migration was removed due to technical reasons")
	return nil
}

func dropVaultAccounts(db migrations.DB) error {
	return nil
}

func init() {
	Collection.MustRegisterTx(loadVaultAccounts, dropVaultAccounts)
}
