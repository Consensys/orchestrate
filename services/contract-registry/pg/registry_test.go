// +build !race

package pg

import (
	"testing"

	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/pg/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/testutils"
)

type PostgresTestSuite struct {
	testutils.ContractRegistryTestSuite
	pg *pgTestUtils.PGTestHelper
}

func (s *PostgresTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
	s.R = NewContractRegistry(s.pg.DB)
}

func (s *PostgresTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
}

func (s *PostgresTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *PostgresTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func TestPostgres(t *testing.T) {
	s := new(PostgresTestSuite)
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	suite.Run(t, s)
}
