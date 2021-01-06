// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres/migrations"
)

const (
	sigHash = "0x12312412"
	codeHash = "0x12324124"
	chainID = "2017"
	address = "0x1234"
)

type eventTestSuite struct {
	suite.Suite
	agents *PGAgents
	pg     *pgTestUtils.PGTestHelper
}

func TestPGEvent(t *testing.T) {
	s := new(eventTestSuite)
	suite.Run(t, s)
}

func (s *eventTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *eventTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *eventTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *eventTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *eventTestSuite) TestPGEvent_InsertMultipleTx() {
	ctx := context.Background()
	s.T().Run("should insert multiple models successfully", func(t *testing.T) {
		events := []*models.EventModel{
			{
				Codehash:          codeHash,
				SigHash:           sigHash,
				IndexedInputCount: 0,
				ABI:               "ABI",
			},
			{
				Codehash:          "codeHash1",
				SigHash:           "sigHash1",
				IndexedInputCount: 1,
				ABI:               "ABI1",
			},
		}
		err := s.agents.Event().InsertMultiple(ctx, events)

		assert.NoError(t, err)
		assert.Equal(t, 1, events[0].ID)
		assert.Equal(t, 2, events[1].ID)
	})

	s.T().Run("should insert multiple models successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx2 := postgres.WithTx(ctx, tx)

		events := []*models.EventModel{
			{
				Codehash:          codeHash,
				SigHash:           sigHash,
				IndexedInputCount: 0,
				ABI:               "ABI",
			},
		}
		err := s.agents.Event().InsertMultiple(ctx2, events)
		_ = tx.Commit()

		assert.NoError(t, err)
		assert.NotEmpty(t, events[0].ID)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		events := []*models.EventModel{
			{
				Codehash:          codeHash,
				SigHash:           sigHash,
				IndexedInputCount: 0,
				ABI:               "ABI",
			},
		}
		err := s.agents.Event().InsertMultiple(ctx, events)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *eventTestSuite) TestPGEvent_FindOneByAccountAndSigHash() {
	ctx := context.Background()

	s.T().Run("should find one event successfully", func(t *testing.T) {
		s.insertEvents(ctx)

		event, err := s.agents.Event().FindOneByAccountAndSigHash(ctx, chainID, address, sigHash, 0)

		assert.NoError(t, err)
		assert.Equal(t, event.ABI, "ABI")
	})

	s.T().Run("should return NotFoundError if no event is found", func(t *testing.T) {
		_, err := s.agents.Event().FindOneByAccountAndSigHash(ctx, "unknown", "unknown", "unknown", 0)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		_, err := s.agents.Event().FindOneByAccountAndSigHash(ctx, chainID, address, sigHash, 0)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *eventTestSuite) TestPGEvent_FindDefaultBySigHash() {
	ctx := context.Background()
	s.T().Run("should find default events successfully", func(t *testing.T) {
		s.insertEvents(ctx)
	
		defaultEvents, err := s.agents.Event().FindDefaultBySigHash(ctx, sigHash, 0)
	
		assert.NoError(t, err)
		assert.Len(t, defaultEvents, 1)
	})

	s.T().Run("should return NotFoundError if no default events are found", func(t *testing.T) {
		defaultEvents, err := s.agents.Event().FindDefaultBySigHash(ctx, "unknown", 0)
		assert.NoError(t, err)
		assert.Empty(t, defaultEvents)
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		s.pg.DropTestDB(t)
	
		_, err := s.agents.Event().FindDefaultBySigHash(ctx, sigHash, 0)
		assert.Error(t, err)
		assert.True(t, errors.IsPostgresConnectionError(err))
	
		s.pg.InitTestDB(t)
	})
}

func (s *eventTestSuite) insertEvents(ctx context.Context) {
	_ = s.agents.CodeHash().Insert(ctx, &models.CodehashModel{
		ChainID:  chainID,
		Address:  address,
		Codehash: codeHash,
	})

	_ = s.agents.CodeHash().Insert(ctx, &models.CodehashModel{
		ChainID:  chainID,
		Address:  "address1",
		Codehash: "codeHash1",
	})

	events := []*models.EventModel{
		{
			Codehash:          codeHash,
			SigHash:           sigHash,
			IndexedInputCount: 0,
			ABI:               "ABI",
		},
		{
			Codehash:          "codeHash1",
			SigHash:           "sigHash1",
			IndexedInputCount: 1,
			ABI:               "ABI1",
		},
	}

	_ = s.agents.Event().InsertMultiple(ctx, events)
}
