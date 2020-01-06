// +build !race

package pg

import (
	"context"
	"testing"

	"github.com/go-pg/pg"
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

func (s *ErrorTestSuite) TestGetNodes() {
	_, err := s.Store.GetNodes(context.Background())
	assert.Error(s.T(), err, "Should get nodes with errors")
}

func (s *ErrorTestSuite) TestGetNodesByTenantID() {
	_, err := s.Store.GetNodesByTenantID(context.Background(), "test")
	assert.Error(s.T(), err, "Should get nodes with errors")
}

func (s *ErrorTestSuite) TestGetNodeByName() {
	_, err := s.Store.GetNodeByName(context.Background(), "test", "test")
	assert.Error(s.T(), err, "Should get node with errors")
}

func (s *ErrorTestSuite) TestGetNodeByID() {
	_, err := s.Store.GetNodeByID(context.Background(), "test")
	assert.Error(s.T(), err, "Should get node with errors")
}

func (s *ErrorTestSuite) TestUpdateNodeByName() {
	err := s.Store.UpdateNodeByName(context.Background(), &types.Node{Name: "test"})
	assert.Error(s.T(), err, "Should update node with errors")
}

func (s *ErrorTestSuite) TestUpdateBlockPositionByName() {
	err := s.Store.UpdateBlockPositionByName(context.Background(), "test", "test", 777)
	assert.Error(s.T(), err, "Should update node with errors")
}

func (s *ErrorTestSuite) TestUpdateNodeByID() {
	err := s.Store.UpdateNodeByID(context.Background(), &types.Node{ID: "test"})
	assert.Error(s.T(), err, "Should update node with errors")
}

func (s *ErrorTestSuite) TestUpdateBlockPositionByID() {
	err := s.Store.UpdateBlockPositionByID(context.Background(), "test", 777)
	assert.Error(s.T(), err, "Should update node with errors")
}

func (s *ErrorTestSuite) TestDeleteNodeByName() {
	err := s.Store.DeleteNodeByName(context.Background(), &types.Node{Name: "test"})
	assert.Error(s.T(), err, "Should update node with errors")
}

func (s *ErrorTestSuite) TestDeleteNodeByID() {
	err := s.Store.DeleteNodeByID(context.Background(), "test")
	assert.Error(s.T(), err, "Should update node with errors")
}
