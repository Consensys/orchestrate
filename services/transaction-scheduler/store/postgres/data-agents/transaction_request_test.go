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

type txRequestTestSuite struct {
	suite.Suite
	dataagent  *PGTransactionRequest
	scheduleDA *PGSchedule
	pg         *pgTestUtils.PGTestHelper
}

func TestPGTransactionRequest(t *testing.T) {
	s := new(txRequestTestSuite)
	suite.Run(t, s)
}

func (s *txRequestTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *txRequestTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.scheduleDA = NewPGSchedule(s.pg.DB)
	s.dataagent = NewPGTransactionRequest(s.pg.DB)
}

func (s *txRequestTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *txRequestTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *txRequestTestSuite) TestPGTransactionRequest_SelectOrInsert() {
	ctx := context.Background()
	schedule := testutils.FakeSchedule()
	_ = s.scheduleDA.Insert(ctx, schedule)

	s.T().Run("should insert model successfully", func(t *testing.T) {
		txRequest := testutils.FakeTxRequest(schedule.ID)
		err := s.dataagent.SelectOrInsert(ctx, txRequest)

		assert.Nil(t, err)
		assert.Equal(t, 1, txRequest.ID)
	})

	s.T().Run("Does nothing if idempotency key is already used and returns request", func(t *testing.T) {
		txRequest0 := testutils.FakeTxRequest(schedule.ID)
		err := s.dataagent.SelectOrInsert(ctx, txRequest0)
		assert.Nil(t, err)

		txRequest1 := testutils.FakeTxRequest(schedule.ID)
		txRequest1.IdempotencyKey = txRequest0.IdempotencyKey
		_ = s.dataagent.SelectOrInsert(ctx, txRequest1)

		assert.Equal(t, txRequest0.IdempotencyKey, txRequest1.IdempotencyKey)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		txRequest := testutils.FakeTxRequest(schedule.ID)
		err := s.dataagent.SelectOrInsert(ctx, txRequest)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *txRequestTestSuite) TestPGTransactionRequest_FindOneByIdempotencyKey() {
	ctx := context.Background()
	schedule := testutils.FakeSchedule()
	_ = s.scheduleDA.Insert(ctx, schedule)

	txRequest := testutils.FakeTxRequest(schedule.ID)
	_ = s.dataagent.SelectOrInsert(ctx, txRequest)

	s.T().Run("should find request successfully", func(t *testing.T) {
		txRequestRetrieved, err := s.dataagent.FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey)

		assert.Nil(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, txRequestRetrieved.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.ID, txRequestRetrieved.Schedule.ID)
	})

	s.T().Run("should return NotFoundError if request is not found", func(t *testing.T) {
		_, err := s.dataagent.FindOneByIdempotencyKey(ctx, "notExisting")
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		txRequestRetrieved, err := s.dataagent.FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey)
		assert.True(t, errors.IsPostgresConnectionError(err))
		assert.Nil(t, txRequestRetrieved)

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}
