package migrations

import (
	"github.com/go-pg/migrations"
)

// Collections holds all migrations
var Collections = migrations.NewCollection()

// Run migrations
func Run(db migrations.DB, a ...string) (oldVersion, newVersion int64, err error) {
	return Collections.Run(db, a...)
}
