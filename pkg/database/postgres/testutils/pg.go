package testutils

import (
	"context"
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
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
	mngr       postgres.Manager
	Collection *migrations.Collection
	opts       *pg.Options
	DB         *pg.DB
}

// NewPGTestHelper creates a new PGTestHelper
func NewPGTestHelper(opts *pg.Options, collection *migrations.Collection) (*PGTestHelper, error) {
	if opts == nil {
		var err error
		opts, err = postgres.NewConfig(viper.GetViper()).PGOptions()
		if err != nil {
			return nil, err
		}
	}

	return &PGTestHelper{
		mngr:       postgres.NewManager(),
		Collection: collection,
		opts:       opts,
	}, nil
}

func (helper *PGTestHelper) Connect(ctx context.Context, opts *pg.Options) *pg.DB {
	return helper.mngr.Connect(ctx, opts)
}

func (helper *PGTestHelper) CreateAndConnect(ctx context.Context, opts *pg.Options) (*pg.DB, error) {
	if opts.Database == "" {
		opts.Database = fmt.Sprintf("test_%s", rand.String(10))
	}

	root := helper.Connect(ctx, helper.opts)
	defer root.Close()
	_, err := root.Exec(`DROP DATABASE IF EXISTS ?;`, pg.SafeQuery(opts.Database))
	if err != nil {
		return nil, err
	}

	_, err = root.Exec(`CREATE DATABASE ?;`, pg.SafeQuery(opts.Database))
	if err != nil {
		return nil, err
	}
	return helper.Connect(ctx, opts), nil
}

func (helper *PGTestHelper) Drop(db *pg.DB) error {
	root := helper.Connect(db.Context(), helper.opts)
	defer root.Close()

	_, err := root.Exec(`DROP DATABASE ?;`, pg.SafeQuery(helper.DB.Options().Database))
	if err != nil {
		return err
	}

	return err
}

func (helper *PGTestHelper) Init(db *pg.DB) error {
	if helper.Collection != nil {
		_, _, err := helper.Collection.Run(db, "init")
		return err
	}
	return nil
}

// Upgrade run migrations 'up'
func (helper *PGTestHelper) Upgrade(db *pg.DB) (oldVersion, newVersion int64, err error) {
	if helper.Collection != nil {
		return helper.Collection.Run(db, "up")
	}
	return
}

// Downgrade run migrations 'reset'
func (helper *PGTestHelper) Downgrade(db *pg.DB) (oldVersion, newVersion int64, err error) {
	if helper.Collection != nil {
		return helper.Collection.Run(db, "reset")
	}
	return
}

// InitTestDB initialize a test database for integration tests
func (helper *PGTestHelper) InitTestDB(t *testing.T) {
	// Create a database for test purpose
	opts := postgres.Copy(helper.opts)
	opts.Database = ""

	db, err := helper.CreateAndConnect(context.Background(), opts)
	require.NoError(t, err, "Creating database should not fail")

	err = helper.Init(db)
	require.NoError(t, err, "Init database should not fail")

	oldVersion, newVersion, err := helper.Upgrade(db)
	require.NoError(t, err, "Upgrading database should not fail")
	t.Logf("Migrated up from version=%v to version=%v", oldVersion, newVersion)

	helper.DB = db
}

func (helper *PGTestHelper) UpgradeTestDB(t *testing.T) {
	oldVersion, newVersion, err := helper.Upgrade(helper.DB)
	require.NoError(t, err, "Upgrade database should not fail")
	t.Logf("Migrated up from version=%v to version=%v", oldVersion, newVersion)
}

func (helper *PGTestHelper) DowngradeTestDB(t *testing.T) {
	oldVersion, newVersion, err := helper.Downgrade(helper.DB)
	require.NoError(t, err, "Downgrade database should not fail")
	t.Logf("Migrated down from version=%v to version=%v", oldVersion, newVersion)
}

// DropTestDB drop test database
func (helper *PGTestHelper) DropTestDB(t *testing.T) {
	helper.DB.Close()
	err := helper.Drop(helper.DB)
	assert.NoError(t, err, "could not drop database")
	helper.DB = nil
}
