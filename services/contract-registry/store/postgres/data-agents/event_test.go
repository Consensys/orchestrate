// +build !race

package dataagents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
)

type eventTestSuite struct {
	suite.Suite
	dataagent  store.EventDataAgent
	codeHashDA store.CodeHashDataAgent
	pg         *pgTestUtils.PGTestHelper
}

func TestPGEvent(t *testing.T) {
	s := new(eventTestSuite)
	suite.Run(t, s)
}

func (s *eventTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *eventTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.codeHashDA = NewPGCodeHash(s.pg.DB)
	s.dataagent = NewPGEvent(s.pg.DB)
}

func (s *eventTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *eventTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *eventTestSuite) TestPGEvent_InsertMultipleTx() {
	s.T().Run("should insert multiple models successfully", func(t *testing.T) {
		events := []*models.EventModel{
			{
				Codehash:          "codeHash",
				SigHash:           "sigHash",
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
		err := s.dataagent.InsertMultiple(context.Background(), &events)

		assert.Nil(t, err)
		assert.Equal(t, 1, events[0].ID)
		assert.Equal(t, 2, events[1].ID)
	})

	s.T().Run("should insert multiple models successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx := postgres.WithTx(context.Background(), tx)

		events := []*models.EventModel{
			{
				Codehash:          "codeHash",
				SigHash:           "sigHash",
				IndexedInputCount: 0,
				ABI:               "ABI",
			},
		}
		err := s.dataagent.InsertMultiple(ctx, &events)
		_ = tx.Commit()

		assert.Nil(t, err)
		assert.NotEmpty(t, events[0].ID)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		events := []*models.EventModel{
			{
				Codehash:          "codeHash",
				SigHash:           "sigHash",
				IndexedInputCount: 0,
				ABI:               "ABI",
			},
		}
		err := s.dataagent.InsertMultiple(context.Background(), &events)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *eventTestSuite) TestPGEvent_FindOneByAccountAndSigHash() {
	s.T().Run("should find one event successfully", func(t *testing.T) {
		s.insertEvents()

		event, err := s.dataagent.FindOneByAccountAndSigHash(context.Background(), &common.AccountInstance{
			ChainId: "chainId",
			Account: "address",
		}, "sigHash", 0)

		assert.Nil(t, err)
		assert.Equal(t, event.ABI, "ABI")
	})

	s.T().Run("should return NotFoundError if no event is found", func(t *testing.T) {
		_, err := s.dataagent.FindOneByAccountAndSigHash(context.Background(), &common.AccountInstance{
			ChainId: "unknown",
			Account: "unknown",
		}, "unknown", 0)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		_, err := s.dataagent.FindOneByAccountAndSigHash(context.Background(), &common.AccountInstance{
			ChainId: "chainId",
			Account: "address",
		}, "sigHash", 0)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *eventTestSuite) TestPGEvent_FindDefaultBySigHash() {
	s.T().Run("should find default events successfully", func(t *testing.T) {
		s.insertEvents()

		defaultEvents, err := s.dataagent.FindDefaultBySigHash(context.Background(), "sigHash", 0)

		assert.Nil(t, err)
		assert.Len(t, defaultEvents, 1)
	})

	s.T().Run("should return NotFoundError if no default events are found", func(t *testing.T) {
		_, err := s.dataagent.FindDefaultBySigHash(context.Background(), "unknown", 0)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		_, err := s.dataagent.FindDefaultBySigHash(context.Background(), "sigHash", 0)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *eventTestSuite) insertEvents() {
	_ = s.codeHashDA.Insert(context.Background(), &models.CodehashModel{
		ChainID:  "chainId",
		Address:  "address",
		Codehash: "codeHash",
	})

	_ = s.codeHashDA.Insert(context.Background(), &models.CodehashModel{
		ChainID:  "chainId",
		Address:  "address1",
		Codehash: "codeHash1",
	})

	events := []*models.EventModel{
		{
			Codehash:          "codeHash",
			SigHash:           "sigHash",
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

	_ = s.dataagent.InsertMultiple(context.Background(), &events)
}
