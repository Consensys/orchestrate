// +build !race

package migrations

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/testutils"
)

type MigrationsTestSuite struct {
	suite.Suite
	pg *testutils.PGTestHelper
}

func (suite *MigrationsTestSuite) SetupSuite() {
	suite.pg.InitTestDB(suite.T())
}

func (suite *MigrationsTestSuite) SetupTest() {
	suite.pg.Upgrade(suite.T())
}

func (suite *MigrationsTestSuite) TearDownTest() {
	suite.pg.Downgrade(suite.T())
}

func (suite *MigrationsTestSuite) TearDownSuite() {
	suite.pg.DropTestDB(suite.T())
}

func (suite *MigrationsTestSuite) TestMigrationVersion() {
	var version int64
	_, err := suite.pg.DB.QueryOne(
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

	n, err := suite.pg.DB.Model().
		Table("pg_catalog.pg_tables").
		Where("tablename = '?'", pg.Q("traces")).
		Count()

	if err != nil {
		suite.T().Errorf("Query failed: %v", err)
	}

	suite.Assert().Equal(1, n, "Trace table should have been created")
}

func (suite *MigrationsTestSuite) TestAddTraceStoreColumns() {
	n, err := suite.pg.DB.Model().
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
	s := new(MigrationsTestSuite)
	s.pg = testutils.NewPGTestHelper(Collection)
	suite.Run(t, s)
}
