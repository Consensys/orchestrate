// +build unit
// +build !race
// +build !integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/suite"
	pgtestutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/postgres/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/testutils"
)

type ModelsTestSuite struct {
	testutils.EnvelopeStoreTestSuite
	pg *pgtestutils.PGTestHelper
}

func (s *ModelsTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
	s.Store = New(s.pg.DB)
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
	s.pg = pgtestutils.NewPGTestHelper(migrations.Collection)
	suite.Run(t, s)
}
