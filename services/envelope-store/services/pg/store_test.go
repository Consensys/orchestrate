// +build !race

package pg

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/envelope-store/services/pg/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/envelope-store/services/testutils"
)

type ModelsTestSuite struct {
	testutils.EnvelopeStoreTestSuite
	pg *testutils.PGTestHelper
}

func (s *ModelsTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
	s.Store = NewEnvelopeStore(s.pg.DB)
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
	s.pg = testutils.NewPGTestHelper(migrations.Collection)
	suite.Run(t, s)
}
