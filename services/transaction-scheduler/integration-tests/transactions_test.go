// +build integration

package integrationtests

import (
	"context"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/testutils"
	"gopkg.in/h2non/gock.v1"
	"testing"
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

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Validation() {
	ctx := context.Background()
	chainUUID := uuid.NewV4().String()

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest()
		txRequest.IdempotencyKey = ""

		resp, err := s.client.SendTransaction(ctx, chainUUID, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail if idempotency key is identical but different params", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Times(2).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendTransactionRequest()

		txResponse, err := s.client.SendTransaction(ctx, chainUUID, txRequest)
		assert.Nil(t, err)

		txRequest.Params.MethodSignature = "differentMethodSignature()"
		txResponse, err = s.client.SendTransaction(ctx, chainUUID, txRequest)
		assert.Nil(t, txResponse)
		assert.True(t, errors.IsConflictedError(err))
	})

	s.T().Run("should fail with 422 if chain does not exist", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(404)
		txRequest := testutils.FakeSendTransactionRequest()

		resp, err := s.client.SendTransaction(ctx, chainUUID, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Transactions() {
	ctx := context.Background()
	chainUUID := uuid.NewV4().String()

	s.T().Run("should send a transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendTransactionRequest()

		txResponse, err := s.client.SendTransaction(ctx, chainUUID, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)
		assert.NotEmpty(t, txResponse.Schedule.UUID)

		scheduleResponse, err := s.client.GetSchedule(ctx, txResponse.Schedule.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, scheduleResponse.Jobs[0].UUID)
		assert.Equal(t, types.StatusStarted, scheduleResponse.Jobs[0].Status)
		assert.Equal(t, txRequest.Params.From, scheduleResponse.Jobs[0].Transaction.From)
		assert.Equal(t, txRequest.Params.To, scheduleResponse.Jobs[0].Transaction.To)

		// TODO: Consume Kafka message and check format
	})
	
	s.T().Run("should send a tessera transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendTesseraRequest()

		txResponse, err := s.client.SendTransaction(ctx, chainUUID, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)
		assert.NotEmpty(t, txResponse.Schedule.UUID)

		scheduleResponse, err := s.client.GetSchedule(ctx, txResponse.Schedule.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, scheduleResponse.Jobs[0].UUID)
		assert.Equal(t, types.StatusStarted, scheduleResponse.Jobs[0].Status)
		assert.Equal(t, txRequest.Params.From, scheduleResponse.Jobs[0].Transaction.From)
		assert.Equal(t, txRequest.Params.To, scheduleResponse.Jobs[0].Transaction.To)
		assert.Equal(t, types.TesseraPrivateTransaction, scheduleResponse.Jobs[0].Type)

		// TODO: Consume Kafka message and check format
	})
	
	s.T().Run("should send a orion transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendOrionRequest()

		txResponse, err := s.client.SendTransaction(ctx, chainUUID, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)
		assert.NotEmpty(t, txResponse.Schedule.UUID)

		scheduleResponse, err := s.client.GetSchedule(ctx, txResponse.Schedule.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, scheduleResponse.Jobs[0].UUID)
		assert.Equal(t, types.StatusStarted, scheduleResponse.Jobs[0].Status)
		assert.Equal(t, txRequest.Params.From, scheduleResponse.Jobs[0].Transaction.From)
		assert.Equal(t, txRequest.Params.To, scheduleResponse.Jobs[0].Transaction.To)
		assert.Equal(t, types.OrionEEATransaction, scheduleResponse.Jobs[0].Type)

		// TODO: Consume Kafka message and check format
	})

	s.T().Run("should succeed if payloads are the same and generate new schedule", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Times(2).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendTransactionRequest()

		txResponse0, err := s.client.SendTransaction(ctx, chainUUID, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		txResponse1, err := s.client.SendTransaction(ctx, chainUUID, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.NotEqual(t, txResponse0.Schedule.UUID, txResponse1.Schedule.UUID)
	})
}
