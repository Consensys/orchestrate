package testutils

import (
	"testing"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/database/postgres"
)

// TODO: all this script should be moved to pkg.git/common
func init() {
	viper.SetDefault("db.user", "postgres")
	_ = viper.BindEnv("db.user", "DB_USER")
	viper.SetDefault("db.password", "postgres")
	_ = viper.BindEnv("db.password", "DB_PASSWORD")
	viper.SetDefault("db.host", "127.0.0.1")
	_ = viper.BindEnv("db.host", "DB_HOST")
	viper.SetDefault("db.port", 5432)
	_ = viper.BindEnv("db.port", "DB_PORT")
	viper.SetDefault("db.database", "postgres")
	_ = viper.BindEnv("db.database", "DB_DATABASE")
}

// PGTestHelper is a suite for integration test of a PostgreSQL database using go-pg
// TODO: move this in pkg.git/common
type PGTestHelper struct {
	Opts       *pg.Options
	DB         *pg.DB
	Collection *migrations.Collection
}

// NewPGTestHelper creates a new PGTestHelper
func NewPGTestHelper(collection *migrations.Collection) *PGTestHelper {
	return &PGTestHelper{
		Opts:       postgres.NewOptions(),
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
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE DATABASE ?;`, pg.Q(testTable))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		log.WithError(err).Warn("could not close Postgres connection")
	}

	helper.DB = pg.Connect(&pg.Options{
		Addr:     helper.Opts.Addr,
		User:     helper.Opts.User,
		Password: helper.Opts.Password,
		Database: "test",
	})
	_, _, err = helper.Collection.Run(helper.DB, "init")
	if err != nil {
		log.Fatal(err)
	}
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
	err := helper.DB.Close()
	if err != nil {
		log.WithError(err).Warn("could not close postgres connection")
	}

	// Drop test Database
	db := pg.Connect(helper.Opts)
	_, err = db.Exec(`DROP DATABASE test;`)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		log.WithError(err).Warn("could not close postgres connection")
	}
}
