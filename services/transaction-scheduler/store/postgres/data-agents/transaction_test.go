// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
)

type txTestSuite struct {
	suite.Suite
	agents *PGAgents
	pg     *pgTestUtils.PGTestHelper
}

func TestPGTx(t *testing.T) {
	s := new(txTestSuite)
	suite.Run(t, s)
}

func (s *txTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *txTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *txTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *txTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *txTestSuite) TestPGTransaction_Insert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully", func(t *testing.T) {
		tx := testutils.FakeTransaction()
		err := s.agents.Transaction().Insert(ctx, tx)

		assert.Nil(t, err)
		assert.NotEmpty(t, tx.ID)
	})

	s.T().Run("should insert private tx model successfully", func(t *testing.T) {
		tx := testutils.FakePrivateTx()
		err := s.agents.Transaction().Insert(ctx, tx)

		assert.Nil(t, err)
		assert.NotEmpty(t, tx.ID)
	})

	s.T().Run("should insert model without UUID successfully", func(t *testing.T) {
		tx := testutils.FakeTransaction()
		tx.UUID = ""
		err := s.agents.Transaction().Insert(ctx, tx)

		assert.Nil(t, err)
		assert.NotEmpty(t, tx.ID)
	})
}

func (s *txTestSuite) TestPGTransaction_Update() {
	ctx := context.Background()

	tx := testutils.FakeTransaction()
	err := s.agents.Transaction().Insert(ctx, tx)
	assert.Nil(s.T(), err)

	s.T().Run("should update model successfully", func(t *testing.T) {
		tx.Hash = "NewHash"
		err := s.agents.Transaction().Update(ctx, tx)

		assert.Nil(t, err)
		assert.NotEmpty(t, tx.Hash, "NewHash")
	})
}

func (s *txTestSuite) TestPGTransaction_ConnectionErr() {
	ctx := context.Background()

	// We drop the DB to make the test fail
	s.pg.DropTestDB(s.T())
	tx := testutils.FakeTransaction()

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		err := s.agents.Transaction().Insert(ctx, tx)
		assert.True(t, errors.IsPostgresConnectionError(err))
	})
	// 
	s.T().Run("should return PostgresConnectionError if update fails", func(t *testing.T) {
		tx.ID = 1
		err := s.agents.Transaction().Update(ctx, tx)
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	// We bring it back up
	s.pg.InitTestDB(s.T())
}
