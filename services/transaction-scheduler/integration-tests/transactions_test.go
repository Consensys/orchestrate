// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/testutils"
	"gopkg.in/h2non/gock.v1"
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
	chainUUID := uuid.Must(uuid.NewV4()).String()

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest()
		txRequest.IdempotencyKey = ""

		resp, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail if idempotency key is identical but different params", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Times(2).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendTransactionRequest()

		txResponse, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)
		assert.Nil(t, err)

		txRequest.Params.MethodSignature = "differentMethodSignature()"
		txResponse, err = s.client.SendContractTransaction(ctx, chainUUID, txRequest)
		assert.Nil(t, txResponse)
		assert.True(t, errors.IsConflictedError(err))
	})

	s.T().Run("should fail with 422 if chain does not exist", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(404)
		txRequest := testutils.FakeSendTransactionRequest()

		resp, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Transactions() {
	ctx := context.Background()
	chainUUID := uuid.Must(uuid.NewV4()).String()

	s.T().Run("should send a transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendTransactionRequest()

		txResponse, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)
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

		txResponse, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)
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

		txResponse, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)
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
	s.T().Run("should send a deploy contract successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeDeployContractRequest()
		txRequest.Params.Args = []string{"123"} // FakeContract arguments

		s.env.contractRegistryResponseFaker.GetContract = func() (*proto.GetContractResponse, error) {
			return &proto.GetContractResponse{
				Contract: testutils2.FakeContract(),
			}, nil
		}
		txResponse, err := s.client.SendDeployTransaction(ctx, chainUUID, txRequest)
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
		assert.Empty(t, scheduleResponse.Jobs[0].Transaction.To)
		assert.Equal(t, types.EthereumTransaction, scheduleResponse.Jobs[0].Type)

		// TODO: Consume Kafka message and check format
	})

	s.T().Run("should send a raw transaction successfully to the transaction sender topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendRawTransactionRequest()

		txResponse, err := s.client.SendRawTransaction(ctx, chainUUID, txRequest)
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
		assert.Equal(t, txRequest.Params.Raw, scheduleResponse.Jobs[0].Transaction.Raw)
		assert.Equal(t, types.EthereumRawTransaction, scheduleResponse.Jobs[0].Type)

		// TODO: Consume Kafka message and check format
	})

	s.T().Run("should succeed if payloads are the same and generate new schedule", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Times(2).Get("/chains/" + chainUUID).Reply(200).JSON(&models.Chain{})
		txRequest := testutils.FakeSendTransactionRequest()

		txResponse0, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		txResponse1, err := s.client.SendContractTransaction(ctx, chainUUID, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.NotEqual(t, txResponse0.Schedule.UUID, txResponse1.Schedule.UUID)
	})
}
