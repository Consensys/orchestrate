// +build integration

package integrationtests

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

// txSchedulerTransactionTestSuite is a test suite for Transaction Scheduler Transactions controller
type txSchedulerTransactionTestSuite struct {
	suite.Suite
	baseURL string
	client  client.TransactionSchedulerClient
	env     *IntegrationEnvironment
}

func (s *txSchedulerTransactionTestSuite) SetupSuite() {
	conf := client.NewConfig(s.baseURL)
	s.client = client.NewHTTPClient(http.NewClient(), conf)
}

/*
func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Validation() {
	ctx := context.Background()

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		txRequest.IdempotencyKey = ""

		resp, err := s.client.SendTransaction(ctx, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail if idempotency key is identical but different params", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()

		txResponse, err := s.client.SendTransaction(ctx, txRequest)
		assert.Nil(t, err)

		txRequest.Params.MethodSignature = "constructor()"
		txResponse, err = s.client.SendTransaction(ctx, txRequest)
		assert.Nil(t, txResponse)
		assert.True(t, errors.IsConflictedError(err))
	})

	s.T().Run("should fail with 422 if chain does not exist", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		txRequest.ChainUUID = uuid.NewV4().String()

		resp, err := s.client.SendTransaction(ctx, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Transactions() {
	ctx := context.Background()

	s.T().Run("should send a transaction successfully to the transaction crafter topic", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedJsonParams, _ := utils.ObjectToJSON(txRequest.Params)
		expectedParams, _ := utils.JSONToMap(expectedJsonParams)

		txResponse, err := s.client.SendTransaction(ctx, txRequest)

		assert.Nil(t, err)

		assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)
		assert.Equal(t, expectedParams, txResponse.Params)
		assert.NotEmpty(t, txResponse.Schedule.UUID)
		assert.Equal(t, txRequest.ChainUUID, txResponse.Schedule.ChainUUID)

		scheduleResponse, err := s.client.GetSchedule(ctx, txResponse.Schedule.UUID)
		assert.Nil(t, err)
		assert.NotEmpty(t, scheduleResponse.Jobs[0].UUID)
		assert.Equal(t, types.JobStatusStarted, scheduleResponse.Jobs[0].Status)
		assert.Equal(t, txRequest.Params.From, scheduleResponse.Jobs[0].Transaction.From)
		assert.Equal(t, txRequest.Params.To, scheduleResponse.Jobs[0].Transaction.To)
	})

	s.T().Run("should succeed if payloads are the same and generate new schedule", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()

		txResponse0, err := s.client.SendTransaction(ctx, txRequest)
		assert.Nil(t, err)
		txResponse1, err := s.client.SendTransaction(ctx, txRequest)
		assert.Nil(t, err)

		assert.NotEqual(t, txResponse0.Schedule.UUID, txResponse1.Schedule.UUID)
	})
}*/
