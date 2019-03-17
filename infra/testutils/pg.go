package testutils

import (
	"testing"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

// PGTestHelper is a suite for integration test of a postgresql database using go-pg
// TODO: move this in pkg.git/common
type PGTestHelper struct {
	Opts       *pg.Options
	DB         *pg.DB
	Collection *migrations.Collection
}

// NewPGTestHelper creates a new PGTestHelper
func NewPGTestHelper(opts *pg.Options, collection *migrations.Collection) *PGTestHelper {
	return &PGTestHelper{
		Opts:       opts,
		Collection: collection,
	}
}

// InitTestDB initialize a test database for integration tests
func (helper *PGTestHelper) InitTestDB(t *testing.T) {
	// Create a test database
	db := pg.Connect(helper.Opts)

	testTable := "test"
	_, err := db.Exec(`DROP DATABASE IF EXISTS ?;`, pg.Q(testTable))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE DATABASE ?;`, pg.Q(testTable))
	if err != nil {
		panic(err)
	}

	db.Close()

	helper.DB = pg.Connect(&pg.Options{
		Addr:     helper.Opts.Addr,
		User:     helper.Opts.User,
		Password: helper.Opts.Password,
		Database: "test",
	})
	helper.Collection.Run(helper.DB, "init")
}

// Upgrade run migrations 'up'
func (helper *PGTestHelper) Upgrade(t *testing.T) {
	oldVersion, newVersion, err := helper.Collection.Run(helper.DB, "up")
	if err != nil {
		t.Errorf("Migrate up: %v\n", err)
	} else {
		t.Logf("Migrated up from version=%v to version=%v\n", oldVersion, newVersion)
	}
}

// Downgrade run migrations 'reset'
func (helper *PGTestHelper) Downgrade(t *testing.T) {
	oldVersion, newVersion, err := helper.Collection.Run(helper.DB, "reset")
	if err != nil {
		t.Errorf("Migrate down: %v\n", err)
	} else {
		t.Logf("Migrated down from version=%v to version=%v\n", oldVersion, newVersion)
	}
}

// DropTestDB drop test database
func (helper *PGTestHelper) DropTestDB(t *testing.T) {
	// Close connection to test database
	helper.DB.Close()

	// Drop test Database
	db := pg.Connect(helper.Opts)
	db.Exec(`DROP DATABASE test;`)
	db.Close()
}
