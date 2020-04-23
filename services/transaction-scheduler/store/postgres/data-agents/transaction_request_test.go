// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
	"testing"
)

type transactionRequestTestSuite struct {
	suite.Suite
	dataagent *PGTransactionRequest
	pg        *pgTestUtils.PGTestHelper
}

func TestPGTransactionRequest(t *testing.T) {
	s := new(transactionRequestTestSuite)
	suite.Run(t, s)
}

func (s *transactionRequestTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *transactionRequestTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.dataagent = NewPGTransactionRequest(s.pg.DB)
}

func (s *transactionRequestTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *transactionRequestTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *transactionRequestTestSuite) TestPGTransactionRequest_Insert() {
	s.T().Run("should insert model successfully", func(t *testing.T) {
		txRequest := testutils.FakeTxRequest()
		err := s.dataagent.Insert(context.Background(), txRequest)

		assert.Nil(t, err)
		assert.Equal(t, txRequest.ID, 1)
	})

	s.T().Run("should insert model successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx := postgres.WithTx(context.Background(), tx)

		txRequest := testutils.FakeTxRequest()
		err := s.dataagent.Insert(ctx, txRequest)
		_ = tx.Commit()

		assert.Nil(t, err)
		assert.Equal(t, txRequest.ID, 2)
	})

	s.T().Run("should return AlreadyExistsError if idempotency key is already used", func(t *testing.T) {
		txRequest0 := testutils.FakeTxRequest()
		err := s.dataagent.Insert(context.Background(), txRequest0)
		assert.Nil(t, err)

		txRequest1 := testutils.FakeTxRequest()
		txRequest1.IdempotencyKey = txRequest0.IdempotencyKey
		err = s.dataagent.Insert(context.Background(), txRequest1)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		err := s.insertTxRequest()
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *transactionRequestTestSuite) insertTxRequest() error {
	txRequest := testutils.FakeTxRequest()
	return s.dataagent.Insert(context.Background(), txRequest)
}
