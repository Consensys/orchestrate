// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
)

type txRequestTestSuite struct {
	suite.Suite
	agents *PGAgents
	pg     *pgTestUtils.PGTestHelper
}

func TestPGTransactionRequest(t *testing.T) {
	s := new(txRequestTestSuite)
	suite.Run(t, s)
}

func (s *txRequestTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *txRequestTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *txRequestTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *txRequestTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *txRequestTestSuite) TestPGTransactionRequest_Insert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully if uuid is not defined", func(t *testing.T) {
		txRequest := testutils.FakeTxRequest(0)
		txRequest.UUID = ""
		err := insertTxRequest(ctx, s.agents, txRequest)

		assert.NoError(t, err)
		assert.NotEmpty(t, txRequest.ID)
		assert.NotEmpty(t, txRequest.UUID)
	})

	s.T().Run("should insert model successfully if uuid is already set", func(t *testing.T) {
		txRequest := testutils.FakeTxRequest(0)
		txRequestUUID := txRequest.UUID

		err := insertTxRequest(ctx, s.agents, txRequest)

		assert.NoError(t, err)
		assert.NotEmpty(t, txRequest.ID)
		assert.Equal(t, txRequestUUID, txRequest.UUID)
	})
}

func (s *txRequestTestSuite) TestPGTransactionRequest_FindOneByIdempotencyKey() {
	ctx := context.Background()
	txRequest := testutils.FakeTxRequest(0)
	err := insertTxRequest(ctx, s.agents, txRequest)
	assert.NoError(s.T(), err)

	s.T().Run("should find request successfully", func(t *testing.T) {
		txRequestRetrieved, err := s.agents.TransactionRequest().
			FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, txRequest.Schedule.TenantID)

		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, txRequestRetrieved.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, txRequestRetrieved.Schedule.UUID)
	})

	s.T().Run("should return NotFoundError if request is not found", func(t *testing.T) {
		_, err := s.agents.TransactionRequest().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, "randomTenant")
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return NotFoundError if request is not found", func(t *testing.T) {
		_, err := s.agents.TransactionRequest().FindOneByIdempotencyKey(ctx, "notExisting", txRequest.Schedule.TenantID)
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *txRequestTestSuite) TestPGTransactionRequest_FindOneByUUID() {
	ctx := context.Background()
	txRequest := testutils.FakeTxRequest(0)
	err := insertTxRequest(ctx, s.agents, txRequest)
	assert.Nil(s.T(), err)

	s.T().Run("should find request successfully for empty tenant", func(t *testing.T) {
		txRequestRetrieved, err := s.agents.TransactionRequest().FindOneByUUID(ctx, txRequest.UUID, []string{multitenancy.Wildcard})

		assert.NoError(t, err)
		assert.Equal(t, txRequest.UUID, txRequestRetrieved.UUID)
		assert.Equal(t, txRequest.Schedule.UUID, txRequestRetrieved.Schedule.UUID)
	})

	s.T().Run("should find request successfully for default tenant", func(t *testing.T) {
		txRequestRetrieved, err := s.agents.TransactionRequest().FindOneByUUID(ctx, txRequest.UUID, []string{multitenancy.DefaultTenant})

		assert.NoError(t, err)
		assert.Equal(t, txRequest.UUID, txRequestRetrieved.UUID)
		assert.Equal(t, txRequest.Schedule.UUID, txRequestRetrieved.Schedule.UUID)
	})

	s.T().Run("should return NotFoundError if uuid is not found", func(t *testing.T) {
		_, err := s.agents.TransactionRequest().FindOneByUUID(ctx, uuid.Must(uuid.NewV4()).String(), []string{txRequest.Schedule.TenantID})
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return NotFoundError if tenant is not found", func(t *testing.T) {
		_, err := s.agents.TransactionRequest().FindOneByUUID(ctx, txRequest.UUID, []string{"notExisting"})
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *txRequestTestSuite) TestPGTransactionRequest_Search() {
	ctx := context.Background()
	txRequest := testutils.FakeTxRequest(0)
	err := insertTxRequest(ctx, s.agents, txRequest)
	assert.Nil(s.T(), err)

	s.T().Run("should find requests successfully", func(t *testing.T) {
		filter := &entities.TransactionFilters{}
		txRequestsRetrieved, err := s.agents.TransactionRequest().Search(ctx, filter, []string{multitenancy.Wildcard})

		assert.NoError(t, err)
		assert.Len(t, txRequestsRetrieved, 1)
		assert.Equal(t, txRequest.UUID, txRequestsRetrieved[0].UUID)
	})

	s.T().Run("should find requests successfully by idempotency keys", func(t *testing.T) {
		filter := &entities.TransactionFilters{
			IdempotencyKeys: []string{txRequest.IdempotencyKey},
		}
		txRequestsRetrieved, err := s.agents.TransactionRequest().Search(ctx, filter, []string{multitenancy.Wildcard})

		assert.NoError(t, err)
		assert.Len(t, txRequestsRetrieved, 1)
		assert.Equal(t, txRequest.UUID, txRequestsRetrieved[0].UUID)
	})

	s.T().Run("should return empty array if nothing found in filter", func(t *testing.T) {
		filter := &entities.TransactionFilters{
			IdempotencyKeys: []string{"notExisting"},
		}

		result, err := s.agents.TransactionRequest().Search(ctx, filter, []string{multitenancy.Wildcard})

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	s.T().Run("should return empty array if tenant is not found", func(t *testing.T) {
		filter := &entities.TransactionFilters{}
		result, err := s.agents.TransactionRequest().Search(ctx, filter, []string{"NotExistingTenant"})

		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func (s *txRequestTestSuite) TestPGTransactionRequest_ConnectionErr() {
	ctx := context.Background()

	// We drop the DB to make the test fail
	s.pg.DropTestDB(s.T())
	txRequest := testutils.FakeTxRequest(0)

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		err := insertTxRequest(ctx, s.agents, txRequest)
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		_, err := s.agents.TransactionRequest().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, txRequest.Schedule.TenantID)
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		_, err := s.agents.TransactionRequest().FindOneByUUID(ctx, txRequest.UUID, []string{multitenancy.Wildcard})
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		_, err := s.agents.TransactionRequest().Search(ctx, &entities.TransactionFilters{}, []string{"tenant"})
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	// We bring it back up
	s.pg.InitTestDB(s.T())
}

func insertTxRequest(ctx context.Context, agents *PGAgents, txReq *models.TransactionRequest) error {
	if err := agents.Schedule().Insert(ctx, txReq.Schedule); err != nil {
		return err
	}

	txReq.ScheduleID = &txReq.Schedule.ID
	if err := agents.TransactionRequest().Insert(ctx, txReq); err != nil {
		return err
	}

	return nil
}
