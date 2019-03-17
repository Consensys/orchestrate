package pg

import (
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/pg/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/testutils"
)

type ModelsTestSuite struct {
	testutils.TraceStoreTestSuite
	pg *testutils.PGTestHelper
}

func (suite *ModelsTestSuite) SetupSuite() {
	suite.pg.InitTestDB(suite.T())
	suite.Store = &TraceStore{db: suite.pg.DB}
}

func (suite *ModelsTestSuite) SetupTest() {
	suite.pg.Upgrade(suite.T())
}

func (suite *ModelsTestSuite) TearDownTest() {
	suite.pg.Downgrade(suite.T())
}

func (suite *ModelsTestSuite) TearDownSuite() {
	suite.pg.DropTestDB(suite.T())
}

func TestModels(t *testing.T) {
	s := new(ModelsTestSuite)
	s.pg = testutils.NewPGTestHelper(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	}, migrations.Collection)
	suite.Run(t, s)
}
