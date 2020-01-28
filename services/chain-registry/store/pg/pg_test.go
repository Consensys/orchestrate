// +build !race

package pg

import (
	"context"
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

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

type ErrorTestSuite struct {
	suite.Suite
	Store types.ChainRegistryStore
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
	s.Store = NewChainRegistry(db)
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (s *ErrorTestSuite) TestGetChains() {
	_, err := s.Store.GetChains(context.Background(), nil)
	assert.Error(s.T(), err, "Should get chains with errors")
}

func (s *ErrorTestSuite) TestGetChainsByTenantID() {
	_, err := s.Store.GetChainsByTenantID(context.Background(), "test", nil)
	assert.Error(s.T(), err, "Should get chains with errors")
}

func (s *ErrorTestSuite) TestGetChainByName() {
	_, err := s.Store.GetChainByTenantIDAndName(context.Background(), "test", "test")
	assert.Error(s.T(), err, "Should get chain with errors")
}

func (s *ErrorTestSuite) TestGetChainByUUID() {
	_, err := s.Store.GetChainByUUID(context.Background(), "test")
	assert.Error(s.T(), err, "Should get chain with errors")
}

func (s *ErrorTestSuite) TestUpdateChainByName() {
	err := s.Store.UpdateChainByName(context.Background(), &types.Chain{Name: "test"})
	assert.Error(s.T(), err, "Should update chain with errors")
}

func (s *ErrorTestSuite) TestUpdateBlockPositionByName() {
	err := s.Store.UpdateBlockPositionByName(context.Background(), "test", "test", 777)
	assert.Error(s.T(), err, "Should update chain with errors")
}

func (s *ErrorTestSuite) TestUpdateChainByUUID() {
	err := s.Store.UpdateChainByUUID(context.Background(), &types.Chain{UUID: "test"})
	assert.Error(s.T(), err, "Should update chain with errors")
}

func (s *ErrorTestSuite) TestUpdateBlockPositionByUUID() {
	err := s.Store.UpdateBlockPositionByUUID(context.Background(), "test", 777)
	assert.Error(s.T(), err, "Should update chain with errors")
}

func (s *ErrorTestSuite) TestDeleteChainByName() {
	err := s.Store.DeleteChainByName(context.Background(), &types.Chain{Name: "test"})
	assert.Error(s.T(), err, "Should update chain with errors")
}

func (s *ErrorTestSuite) TestDeleteChainByUUID() {
	err := s.Store.DeleteChainByUUID(context.Background(), "test")
	assert.Error(s.T(), err, "Should update chain with errors")
}
