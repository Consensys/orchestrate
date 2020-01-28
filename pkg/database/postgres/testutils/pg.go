package testutils

import (
	"fmt"
	"testing"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"k8s.io/apimachinery/pkg/util/rand"
)

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
type PGTestHelper struct {
	Opts       *pg.Options
	DB         *pg.DB
	Collection *migrations.Collection
	TestDBName string
}

// NewPGTestHelper creates a new PGTestHelper
func NewPGTestHelper(collection *migrations.Collection) *PGTestHelper {
	return &PGTestHelper{
		Opts:       postgres.NewOptions(),
		Collection: collection,
		TestDBName: fmt.Sprintf("test_%s", rand.String(10)),
	}
}

// InitTestDB initialize a test database for integration tests
func (helper *PGTestHelper) InitTestDB(t *testing.T) {
	db := pg.Connect(helper.Opts)
	_, err := db.Exec(`DROP DATABASE IF EXISTS ?;`, pg.Q(helper.TestDBName))
	assert.NoError(t, err, "could not drop database")

	_, err = db.Exec(`CREATE DATABASE ?;`, pg.Q(helper.TestDBName))
	assert.NoError(t, err, "could not create database")

	err = db.Close()
	if err != nil {
		log.WithError(err).Warn("could not close Postgres connection")
	}

	helper.DB = pg.Connect(&pg.Options{
		Addr:     helper.Opts.Addr,
		User:     helper.Opts.User,
		Password: helper.Opts.Password,
		Database: helper.TestDBName,
	})
	_, _, err = helper.Collection.Run(helper.DB, "init")
	if err != nil {
		log.WithError(err).Fatal("could not init database")
	}
}

// Upgrade run migrations 'up'
func (helper *PGTestHelper) Upgrade(t *testing.T) {
	oldVersion, newVersion, err := helper.Collection.Run(helper.DB, "up")
	assert.NoError(t, err, "Failed migrate up")

	t.Logf("Migrated up from version=%v to version=%v", oldVersion, newVersion)
}

// Downgrade run migrations 'reset'
func (helper *PGTestHelper) Downgrade(t *testing.T) {
	oldVersion, newVersion, err := helper.Collection.Run(helper.DB, "reset")
	assert.NoError(t, err, "Failed migrate down")

	t.Logf("Migrated down from version=%v to version=%v", oldVersion, newVersion)
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
	_, err = db.Exec(`DROP DATABASE ?;`, pg.Q(helper.TestDBName))
	assert.NoError(t, err, "could not drop database")

	err = db.Close()
	if err != nil {
		log.WithError(err).Warn("could not close postgres connection")
	}
}
