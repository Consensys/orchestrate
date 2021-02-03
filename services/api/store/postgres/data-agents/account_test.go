// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres/migrations"
)

type accountTestSuite struct {
	suite.Suite
	agents   *PGAgents
	pg       *pgTestUtils.PGTestHelper
	tenantID string
}

func TestPGAccount(t *testing.T) {
	s := new(accountTestSuite)
	suite.Run(t, s)
}

func (s *accountTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.tenantID = "tenantID"
	s.pg.InitTestDB(s.T())
}

func (s *accountTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *accountTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *accountTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *accountTestSuite) TestPGAccount_Insert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), acc.ID)
	})

	s.T().Run("should fail to insert same alias twice", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)

		acc.Address = ethcommon.HexToAddress("0x322").String()
		err = s.agents.Account().Insert(ctx, acc)
		assert.Error(s.T(), err)
	})

	s.T().Run("should fail to insert same address twice with same tenant", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)

		acc.Alias = utils.RandString(10)
		err = s.agents.Account().Insert(ctx, acc)
		assert.Error(s.T(), err)
	})
}

func (s *accountTestSuite) TestPGAccount_FindOneByAddress() {
	ctx := context.Background()

	s.T().Run("should find one model by address successfully", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)

		iden2, err := s.agents.Account().FindOneByAddress(ctx, acc.Address, []string{acc.TenantID})
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), acc, iden2)
	})

	s.T().Run("should fail to find mode with different tenant", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)

		_, err = s.agents.Account().FindOneByAddress(ctx, acc.Address, []string{"Not tenant"})
		assert.Error(s.T(), err)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *accountTestSuite) TestPGAccount_Search() {
	ctx := context.Background()

	s.T().Run("should search model successfully", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)

		accs, err := s.agents.Account().Search(ctx, &entities.AccountFilters{Aliases: []string{acc.Alias}}, []string{acc.TenantID})
		assert.NoError(s.T(), err)
		assert.Len(s.T(), accs, 1)

		// Re-write insert updated fields to validate remaining properties
		accs[0].ID = acc.ID
		accs[0].CreatedAt = acc.CreatedAt
		accs[0].UpdatedAt = acc.UpdatedAt
		assert.Equal(s.T(), accs[0], acc)
	})

	s.T().Run("should search model successfully, filter by tenant", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)

		accs, err := s.agents.Account().Search(ctx, &entities.AccountFilters{Aliases: []string{acc.Alias}}, []string{"invalidTenant"})
		assert.NoError(s.T(), err)
		assert.Len(s.T(), accs, 0)
	})
}

func (s *accountTestSuite) TestPGAccount_Update() {
	ctx := context.Background()

	s.T().Run("should update model successfully", func(t *testing.T) {
		acc := testutils2.FakeAccountModel()
		err := s.agents.Account().Insert(ctx, acc)
		assert.NoError(s.T(), err)

		acc.Attributes = map[string]string{
			"newAttr3": "newVal3",
		}
		acc.Alias = "NewAlias"
		err = s.agents.Account().Update(ctx, acc)
		assert.NoError(s.T(), err)

		iden2, err := s.agents.Account().FindOneByAddress(ctx, acc.Address, []string{acc.TenantID})
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), acc.Alias, iden2.Alias)
		assert.Equal(s.T(), acc.Attributes, iden2.Attributes)
	})
}
