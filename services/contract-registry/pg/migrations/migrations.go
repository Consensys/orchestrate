package migrations

import (
	"github.com/go-pg/migrations"
)

// Collection holds all migrations
var Collection = migrations.NewCollection()

// Forces disablement of SQLAutodiscover before any init() function can append
var _ = func() error { Collection.DisableSQLAutodiscover(true); return nil }() //nolint:unparam // reason

// Run migrations
func Run(db migrations.DB, a ...string) (oldVersion, newVersion int64, err error) {
	return Collection.Run(db, a...)
}
