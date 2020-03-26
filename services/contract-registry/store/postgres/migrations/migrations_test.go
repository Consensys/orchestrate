// +build !race

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
	s.pg.Upgrade(s.T())
}

func (s *MigrationsTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
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
	s.Assert().Equal(int64(2), version, "Migration should be on correct version")
}

func (s *MigrationsTestSuite) TestCreateArtifactsTable() {
	n, err := s.pg.DB.Model().
		Table("pg_catalog.pg_tables").
		Where("tablename = '?'", pg.SafeQuery("codehashes")).
		Count()

	s.Assert().NoError(err, "Query failed")
	s.Assert().Equal(1, n, "Table should have been created")
}

func (s *MigrationsTestSuite) TestAddEnvelopeStoreColumns() {
}

func TestMigrations(t *testing.T) {
	s := new(MigrationsTestSuite)
	s.pg = testutils.NewPGTestHelper(Collection)
	suite.Run(t, s)
}
