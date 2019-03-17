package migrations

import (
	"github.com/go-pg/migrations"
)

// Collection holds all migrations
var Collection = migrations.NewCollection()

// Run migrations
func Run(db migrations.DB, a ...string) (oldVersion, newVersion int64, err error) {
	return Collection.Run(db, a...)
}
