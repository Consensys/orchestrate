// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"
	"github.com/stretchr/testify/assert"

	pgTestUtils "github.com/consensys/orchestrate/pkg/toolkit/database/postgres/testutils"
	"github.com/consensys/orchestrate/services/api/store/postgres/migrations"
	"github.com/stretchr/testify/suite"
)

type chainTestSuite struct {
	suite.Suite
	agents         *PGAgents
	pg             *pgTestUtils.PGTestHelper
	allowedTenants []string
	tenantID       string
	username       string
}

func TestPGChain(t *testing.T) {
	s := new(chainTestSuite)
	suite.Run(t, s)
}

func (s *chainTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.tenantID = "tenantID"
	s.allowedTenants = []string{s.tenantID, "_"}
	s.username = "username"
	s.pg.InitTestDB(s.T())
}

func (s *chainTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *chainTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *chainTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *chainTestSuite) TestInsert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully", func(t *testing.T) {
		chain := testutils.FakeChainModel()
		err := s.agents.Chain().Insert(ctx, chain)
		assert.NoError(s.T(), err)

		assert.NoError(t, err)
		assert.NotEmpty(t, chain.UUID)
	})

	s.T().Run("should insert model without UUID successfully", func(t *testing.T) {
		chain := testutils.FakeChainModel()
		chain.UUID = ""
		chain.Name = "insert_2"
		err := s.agents.Chain().Insert(ctx, chain)
		assert.NoError(s.T(), err)

		assert.NoError(t, err)
		assert.NotEmpty(t, chain.UUID)
	})
}

func (s *chainTestSuite) TestUpdate() {
	ctx := context.Background()
	chain := testutils.FakeChainModel()
	err := s.agents.Chain().Insert(ctx, chain)
	assert.NoError(s.T(), err)

	s.T().Run("should update model successfully", func(t *testing.T) {
		newChain := testutils.FakeChainModel()
		newChain.OwnerID = s.username
		newChain.TenantID = s.tenantID
		newChain.ListenerCurrentBlock = 666
		newChain.UUID = chain.UUID

		err = s.agents.Chain().Update(ctx, newChain, s.allowedTenants, s.username)
		assert.NoError(t, err)

		chainRetrieved, _ := s.agents.Chain().FindOneByUUID(ctx, newChain.UUID, s.allowedTenants, s.username)
		assert.Equal(t, newChain.ListenerCurrentBlock, chainRetrieved.ListenerCurrentBlock)
	})

	s.T().Run("should update model successfully with tenant", func(t *testing.T) {
		newChain := testutils.FakeChainModel()
		newChain.OwnerID = s.username
		newChain.TenantID = s.tenantID
		newChain.ListenerCurrentBlock = 666
		newChain.UUID = chain.UUID

		err = s.agents.Chain().Update(ctx, newChain, s.allowedTenants, s.username)

		assert.NoError(t, err)

		chainRetrieved, _ := s.agents.Chain().FindOneByUUID(ctx, newChain.UUID, s.allowedTenants, s.username)
		assert.Equal(t, newChain.ListenerCurrentBlock, chainRetrieved.ListenerCurrentBlock)
	})
}

func (s *chainTestSuite) TestFindOneByUUID() {
	ctx := context.Background()
	chain := testutils.FakeChainModel()
	chain.OwnerID = s.username
	chain.TenantID = s.tenantID
	err := s.agents.Chain().Insert(ctx, chain)
	assert.NoError(s.T(), err)

	s.T().Run("should get model successfully", func(t *testing.T) {
		chainRetrieved, err := s.agents.Chain().FindOneByUUID(ctx, chain.UUID, s.allowedTenants, s.username)

		assert.NoError(t, err)
		assert.Equal(t, chain.UUID, chainRetrieved.UUID)
	})

	s.T().Run("should get model successfully as tenant", func(t *testing.T) {
		chainRetrieved, err := s.agents.Chain().FindOneByUUID(ctx, chain.UUID, s.allowedTenants, s.username)

		assert.NoError(t, err)
		assert.NotEmpty(t, chainRetrieved.UUID)
	})

	s.T().Run("should return NotFoundError if select fails", func(t *testing.T) {
		_, err := s.agents.Chain().FindOneByUUID(ctx, "b6fe7a2a-1a4d-49ca-99d8-8a34aa495ef0", s.allowedTenants, s.username)
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *chainTestSuite) TestSearch() {
	ctx := context.Background()

	chain0 := testutils.FakeChainModel()
	chain0.Name = "chain0"
	chain0.OwnerID = s.username
	chain0.TenantID = s.tenantID
	err := s.agents.Chain().Insert(ctx, chain0)
	assert.NoError(s.T(), err)

	chain1 := testutils.FakeChainModel()
	chain1.Name = "chain1"
	chain1.OwnerID = s.username
	chain1.TenantID = s.tenantID
	err = s.agents.Chain().Insert(ctx, chain1)
	assert.NoError(s.T(), err)

	s.T().Run("should find models successfully with filters", func(t *testing.T) {
		filters := &entities.ChainFilters{
			Names: []string{chain0.Name},
		}

		retrievedChains, err := s.agents.Chain().Search(ctx, filters, s.allowedTenants, s.username)

		assert.NoError(t, err)
		assert.Equal(t, chain0.UUID, retrievedChains[0].UUID)
		assert.Len(t, retrievedChains, 1)
	})

	s.T().Run("should not find any model by names", func(t *testing.T) {
		filters := &entities.ChainFilters{
			Names: []string{"0x3"},
		}

		retrievedChains, err := s.agents.Chain().Search(ctx, filters, s.allowedTenants, s.username)

		assert.NoError(t, err)
		assert.Empty(t, retrievedChains)
	})

	s.T().Run("should find every inserted model successfully", func(t *testing.T) {
		filters := &entities.ChainFilters{}
		retrievedChains, err := s.agents.Chain().Search(ctx, filters, s.allowedTenants, s.username)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(retrievedChains))
	})
}

func (s *chainTestSuite) TestConnectionErr() {
	ctx := context.Background()

	// We drop the DB to make the test fail
	s.pg.DropTestDB(s.T())
	chain := testutils.FakeChainModel()
	chain.OwnerID = s.username
	chain.TenantID = s.tenantID
	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		err := s.agents.Chain().Insert(ctx, chain)
		assert.True(t, errors.IsInternalError(err))
	})

	s.T().Run("should return PostgresConnectionError if update fails", func(t *testing.T) {
		err := s.agents.Chain().Update(ctx, chain, s.allowedTenants, s.username)
		assert.True(t, errors.IsInternalError(err))
	})

	s.T().Run("should return PostgresConnectionError if findOne fails", func(t *testing.T) {
		_, err := s.agents.Chain().FindOneByUUID(ctx, chain.UUID, s.allowedTenants, s.username)
		assert.True(t, errors.IsInternalError(err))
	})

	// We bring it back up
	s.pg.InitTestDB(s.T())
}
