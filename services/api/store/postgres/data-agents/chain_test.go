// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/store/models/testutils"
	"testing"

	"github.com/stretchr/testify/suite"
	pgTestUtils "github.com/ConsenSys/orchestrate/pkg/database/postgres/testutils"
	"github.com/ConsenSys/orchestrate/services/api/store/postgres/migrations"
)

type chainTestSuite struct {
	suite.Suite
	agents   *PGAgents
	pg       *pgTestUtils.PGTestHelper
	tenantID string
}

func TestPGChain(t *testing.T) {
	s := new(chainTestSuite)
	suite.Run(t, s)
}

func (s *chainTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.tenantID = "tenantID"
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
		newChain.ListenerCurrentBlock = 666
		newChain.UUID = chain.UUID

		err = s.agents.Chain().Update(ctx, newChain, []string{multitenancy.Wildcard})
		assert.NoError(t, err)

		chainRetrieved, _ := s.agents.Chain().FindOneByUUID(ctx, newChain.UUID, []string{multitenancy.Wildcard})
		assert.Equal(t, newChain.ListenerCurrentBlock, chainRetrieved.ListenerCurrentBlock)
	})

	s.T().Run("should update model successfully with tenant", func(t *testing.T) {
		newChain := testutils.FakeChainModel()
		newChain.ListenerCurrentBlock = 666
		newChain.UUID = chain.UUID

		err = s.agents.Chain().Update(ctx, newChain, []string{newChain.TenantID})

		assert.NoError(t, err)

		chainRetrieved, _ := s.agents.Chain().FindOneByUUID(ctx, newChain.UUID, []string{multitenancy.Wildcard})
		assert.Equal(t, newChain.ListenerCurrentBlock, chainRetrieved.ListenerCurrentBlock)
	})
}

func (s *chainTestSuite) TestFindOneByUUID() {
	ctx := context.Background()
	chain := testutils.FakeChainModel()
	err := s.agents.Chain().Insert(ctx, chain)
	assert.NoError(s.T(), err)

	s.T().Run("should get model successfully", func(t *testing.T) {
		chainRetrieved, err := s.agents.Chain().FindOneByUUID(ctx, chain.UUID, []string{multitenancy.Wildcard})

		assert.NoError(t, err)
		assert.Equal(t, chain.UUID, chainRetrieved.UUID)
	})

	s.T().Run("should get model successfully as tenant", func(t *testing.T) {
		chainRetrieved, err := s.agents.Chain().FindOneByUUID(ctx, chain.UUID, []string{s.tenantID})

		assert.NoError(t, err)
		assert.NotEmpty(t, chainRetrieved.UUID)
	})

	s.T().Run("should return NotFoundError if select fails", func(t *testing.T) {
		_, err := s.agents.Chain().FindOneByUUID(ctx, "b6fe7a2a-1a4d-49ca-99d8-8a34aa495ef0", []string{s.tenantID})
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *chainTestSuite) TestSearch() {
	ctx := context.Background()

	chain0 := testutils.FakeChainModel()
	chain0.Name = "chain0"
	err := s.agents.Chain().Insert(ctx, chain0)
	assert.NoError(s.T(), err)

	chain1 := testutils.FakeChainModel()
	chain1.Name = "chain1"
	err = s.agents.Chain().Insert(ctx, chain1)
	assert.NoError(s.T(), err)

	s.T().Run("should find models successfully with filters", func(t *testing.T) {
		filters := &entities.ChainFilters{
			Names: []string{chain0.Name},
		}

		retrievedChains, err := s.agents.Chain().Search(ctx, filters, []string{s.tenantID})

		assert.NoError(t, err)
		assert.Equal(t, chain0.UUID, retrievedChains[0].UUID)
		assert.Len(t, retrievedChains, 1)
	})

	s.T().Run("should not find any model by names", func(t *testing.T) {
		filters := &entities.ChainFilters{
			Names: []string{"0x3"},
		}

		retrievedChains, err := s.agents.Chain().Search(ctx, filters, []string{s.tenantID})

		assert.NoError(t, err)
		assert.Empty(t, retrievedChains)
	})

	s.T().Run("should find every inserted model successfully", func(t *testing.T) {
		filters := &entities.ChainFilters{}
		retrievedChains, err := s.agents.Chain().Search(ctx, filters, []string{s.tenantID})

		assert.NoError(t, err)
		assert.Equal(t, 2, len(retrievedChains))
	})
}

func (s *chainTestSuite) TestConnectionErr() {
	ctx := context.Background()

	// We drop the DB to make the test fail
	s.pg.DropTestDB(s.T())
	chain := testutils.FakeChainModel()
	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		err := s.agents.Chain().Insert(ctx, chain)
		assert.True(t, errors.IsInternalError(err))
	})

	s.T().Run("should return PostgresConnectionError if update fails", func(t *testing.T) {
		err := s.agents.Chain().Update(ctx, chain, []string{chain.TenantID})
		assert.True(t, errors.IsInternalError(err))
	})

	s.T().Run("should return PostgresConnectionError if findOne fails", func(t *testing.T) {
		_, err := s.agents.Chain().FindOneByUUID(ctx, chain.UUID, []string{chain.TenantID})
		assert.True(t, errors.IsInternalError(err))
	})

	// We bring it back up
	s.pg.InitTestDB(s.T())
}
