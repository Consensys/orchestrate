// +build unit
// +build !race
// +build !integration

package pg_test

import (
	"context"
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

func (s *ChainModelsTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
	s.Store = pgstore.New(s.pg.DB)
}

func (s *ChainModelsTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
}

func (s *ChainModelsTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *ChainModelsTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *FaucetModelsTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
	s.Store = pgstore.New(s.pg.DB)
}

func (s *FaucetModelsTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
}

func (s *FaucetModelsTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *FaucetModelsTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

type ChainModelsTestSuite struct {
	pg *pgTestUtils.PGTestHelper
	testutils.ChainTestSuite
}

func TestModelsChain(t *testing.T) {
	s := new(ChainModelsTestSuite)
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	suite.Run(t, s)
}

type FaucetModelsTestSuite struct {
	pg *pgTestUtils.PGTestHelper
	testutils.FaucetTestSuite
}

func TestModelsFaucet(t *testing.T) {
	s := new(FaucetModelsTestSuite)
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	suite.Run(t, s)
}

type ErrorTestSuite struct {
	suite.Suite
	Store *pgstore.PG
}

func (s *ErrorTestSuite) SetupSuite() {
	options := &pg.Options{
		Addr:     "error",
		User:     "error",
		Password: "error",
		Database: "error",
		PoolSize: 1,
	}
	db := pg.Connect(options)
	s.Store = pgstore.New(db)
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (s *ErrorTestSuite) TestGetChains() {
	_, err := s.Store.GetChains(context.Background(), nil)
	assert.Error(s.T(), err, "Should get chains with errors")
}

func (s *ErrorTestSuite) TestGetChainByUUID() {
	_, err := s.Store.GetChainByUUID(context.Background(), "test")
	assert.Error(s.T(), err, "Should get chain with errors")
}

func (s *ErrorTestSuite) TestUpdateChainByName() {
	err := s.Store.UpdateChainByName(context.Background(), &types.Chain{Name: "test"})
	assert.Error(s.T(), err, "Should update chain with errors")
}

func (s *ErrorTestSuite) TestUpdateChainByUUID() {
	err := s.Store.UpdateChainByUUID(context.Background(), &types.Chain{UUID: "test"})
	assert.Error(s.T(), err, "Should update chain with errors")
}

func (s *ErrorTestSuite) TestDeleteChainByUUID() {
	err := s.Store.DeleteChainByUUID(context.Background(), "test")
	assert.Error(s.T(), err, "Should update chain with errors")
}
