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
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/postgres/migrations"
)

type identityTestSuite struct {
	suite.Suite
	agents   *PGAgents
	pg       *pgTestUtils.PGTestHelper
	tenantID string
}

func TestPGJob(t *testing.T) {
	s := new(identityTestSuite)
	suite.Run(t, s)
}

func (s *identityTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.tenantID = "tenantID"
	s.pg.InitTestDB(s.T())
}

func (s *identityTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *identityTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *identityTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *identityTestSuite) TestPGJob_Insert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully", func(t *testing.T) {
		iden := testutils2.FakeIdentityModel()
		err := s.agents.Identity().Insert(ctx, iden)
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), iden.ID)
	})

	s.T().Run("should fail to insert same alias twice", func(t *testing.T) {
		iden := testutils2.FakeIdentityModel()
		err := s.agents.Identity().Insert(ctx, iden)
		assert.NoError(s.T(), err)

		iden.Address = ethcommon.HexToAddress("0x322").String()
		err = s.agents.Identity().Insert(ctx, iden)
		assert.Error(s.T(), err)
	})
}

func (s *identityTestSuite) TestPGJob_Search() {
	ctx := context.Background()

	s.T().Run("should search model successfully", func(t *testing.T) {
		iden := testutils2.FakeIdentityModel()
		err := s.agents.Identity().Insert(ctx, iden)
		assert.NoError(s.T(), err)

		idens, err := s.agents.Identity().Search(ctx, &entities.IdentityFilters{Aliases:[]string{iden.Alias}}, []string{iden.TenantID})
		assert.NoError(s.T(), err)
		assert.Len(s.T(), idens, 1)
		
		// Re-write insert updated fields to validate remaining properties
		idens[0].ID = iden.ID
		idens[0].CreatedAt = iden.CreatedAt
		idens[0].UpdatedAt = iden.UpdatedAt
		assert.Equal(s.T(), idens[0], iden)
	})
	
	s.T().Run("should search model successfully, filter by tenant", func(t *testing.T) {
		iden := testutils2.FakeIdentityModel()
		err := s.agents.Identity().Insert(ctx, iden)
		assert.NoError(s.T(), err)

		idens, err := s.agents.Identity().Search(ctx, &entities.IdentityFilters{Aliases:[]string{iden.Alias}}, []string{"invalidTenant"})
		assert.NoError(s.T(), err)
		assert.Len(s.T(), idens, 0)
	})
}
