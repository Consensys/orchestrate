// +build !race

package pg

import (
	"testing"

	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/testutils"
)

type ModelsTestSuite struct {
	testutils.ChainRegistryTestSuite
	pg *pgTestUtils.PGTestHelper
}

func (s *ModelsTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
	s.Store = NewChainRegistry(s.pg.DB)
}

func (s *ModelsTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
}

func (s *ModelsTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *ModelsTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func TestModels(t *testing.T) {
	s := new(ModelsTestSuite)
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	suite.Run(t, s)
}
