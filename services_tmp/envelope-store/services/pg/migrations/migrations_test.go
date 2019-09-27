// +build !race

package migrations

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/services/testutils"
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
		pg.Q("gopg_migrations"),
	)

	if err != nil {
		s.T().Errorf("Error querying version: %v", err)
	}

	expected := int64(3)
	s.Assert().Equal(expected, version, fmt.Sprintf("Migration should be on version=%v", expected))
}

func (s *MigrationsTestSuite) TestCreateEnvelopeTable() {

	n, err := s.pg.DB.Model().
		Table("pg_catalog.pg_tables").
		Where("tablename = '?'", pg.Q("envelopes")).
		Count()

	if err != nil {
		s.T().Errorf("Query failed: %v", err)
	}

	s.Assert().Equal(1, n, "Envelope table should have been created")
}

func (s *MigrationsTestSuite) TestAddEnvelopeStoreColumns() {
	n, err := s.pg.DB.Model().
		Table("information_schema.columns").
		Where("table_name = '?'", pg.Q("envelopes")).
		Count()

	if err != nil {
		s.T().Errorf("Query failed: %v", err)
	}

	expected := 10
	s.Assert().Equal(expected, n, fmt.Sprintf("Envelope table should have %v columns", expected))
}

func TestMigrations(t *testing.T) {
	s := new(MigrationsTestSuite)
	s.pg = testutils.NewPGTestHelper(Collection)
	suite.Run(t, s)
}
