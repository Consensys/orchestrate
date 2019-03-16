package migrations

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/suite"
)

type TraceStoreTestSuite struct {
	suite.Suite
	db *pg.DB
}

func (suite *TraceStoreTestSuite) SetupSuite() {
	// Create a test database
	db := pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})

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

	// Create a connection to test database
	suite.db = pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "test",
	})
	Run(suite.db, "init")
}

func (suite *TraceStoreTestSuite) SetupTest() {
	oldVersion, newVersion, err := Run(suite.db, "up")
	if err != nil {
		suite.T().Errorf("Migrate up: %v", err)
	} else {
		suite.T().Logf("Migrated up from version=%v to version=%v", oldVersion, newVersion)
	}
}

func (suite *TraceStoreTestSuite) TearDownTest() {
	oldVersion, newVersion, err := Run(suite.db, "reset")
	if err != nil {
		suite.T().Errorf("Migrate down: %v", err)
	} else {
		suite.T().Logf("Migrated down from version=%v to version=%v", oldVersion, newVersion)
	}
}

func (suite *TraceStoreTestSuite) TearDownSuite() {
	// Close connection to test database
	suite.db.Close()

	// Drop test Database
	db := pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})
	db.Exec(`DROP DATABASE test;`)
	db.Close()
}

type MigrationsTestSuite struct {
	TraceStoreTestSuite
}

func (suite *MigrationsTestSuite) TestMigrationVersion() {
	var version int64
	_, err := suite.db.QueryOne(
		pg.Scan(&version),
		`SELECT version FROM ? ORDER BY id DESC LIMIT 1`,
		pg.Q("gopg_migrations"),
	)

	if err != nil {
		suite.T().Errorf("Error querying version: %v", err)
	}

	expected := int64(2)
	suite.Assert().Equal(expected, version, fmt.Sprintf("Migration should be on version=%v", expected))
}

func (suite *MigrationsTestSuite) TestCreateTraceTable() {

	n, err := suite.db.Model().
		Table("pg_catalog.pg_tables").
		Where("tablename = '?'", pg.Q("traces")).
		Count()

	if err != nil {
		suite.T().Errorf("Query failed: %v", err)
	}

	suite.Assert().Equal(1, n, "Trace table should have been created")
}

func (suite *MigrationsTestSuite) TestAddTraceStoreColumns() {
	n, err := suite.db.Model().
		Table("information_schema.columns").
		Where("table_name = '?'", pg.Q("traces")).
		Count()

	if err != nil {
		suite.T().Errorf("Query failed: %v", err)
	}

	expected := 10
	suite.Assert().Equal(expected, n, fmt.Sprintf("Trace table should have %v columns", expected))
}

func TestMigrations(t *testing.T) {
	suite.Run(t, new(MigrationsTestSuite))
}
