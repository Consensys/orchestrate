// +build unit
// +build !race
// +build !integration

package migrations

import (
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
)

type MigrationsTestSuite struct {
	suite.Suite
	pg *testutils.PGTestHelper
}

func (s *MigrationsTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
}

func (s *MigrationsTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
}

func (s *MigrationsTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *MigrationsTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *MigrationsTestSuite) TestMigrationVersion() {
	var version int64
	_, err := s.pg.DB.QueryOne(
		pg.Scan(&version),
		`SELECT version FROM ? ORDER BY id DESC LIMIT 1`,
		pg.SafeQuery("gopg_migrations"),
	)

	s.Assert().NoError(err, "Error querying version")
	s.Assert().Equal(int64(3), version, "Migration should be on correct version")
}

func (s *MigrationsTestSuite) TestCreateRequestsTable() {
	n, err := s.pg.DB.Model().
		Table("pg_catalog.pg_tables").
		Where("tablename = '?'", pg.SafeQuery("transaction_requests")).
		Count()

	s.Assert().NoError(err, "Query failed")
	s.Assert().Equal(1, n, "Table should have been created")
}

func (s *MigrationsTestSuite) TestAddRequestsColumns() {
	n, err := s.pg.DB.Model().
		Table("information_schema.columns").
		Where("table_name = '?'", pg.SafeQuery("transaction_requests")).
		Count()

	s.Assert().NoError(err, "Query failed")
	s.Assert().Equal(8, n, "Requests table should have correct number of columns")
}

func TestMigrations(t *testing.T) {
	s := new(MigrationsTestSuite)
	s.pg, _ = testutils.NewPGTestHelper(nil, Collection)
	suite.Run(t, s)
}
