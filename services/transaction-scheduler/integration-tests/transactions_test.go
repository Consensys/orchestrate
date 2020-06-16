// +build integration

package integrationtests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/testutils"
	"gopkg.in/h2non/gock.v1"
)

const (
	waitForEnvelopeTimeOut = 2 * time.Second
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
	chain := testutils2.FakeChain()
	chainModel := &models.Chain{
		Name: chain.Name,
		UUID: chain.UUID,
		TenantID: chain.TenantID,
	}

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest(chain.Name)
		txRequest.IdempotencyKey = ""

		resp, err := s.client.SendContractTransaction(ctx, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail if idempotency key is identical but different params", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest(chain.Name)

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txResponse, err := s.client.SendContractTransaction(ctx, txRequest)
		assert.Nil(t, err)

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txRequest.Params.MethodSignature = "differentMethodSignature()"
		txResponse, err = s.client.SendContractTransaction(ctx, txRequest)
		assert.Nil(t, txResponse)
		assert.True(t, errors.IsConflictedError(err))
	})

	s.T().Run("should fail with 422 if chain does not exist", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(404)
		txRequest := testutils.FakeSendTransactionRequest(chain.Name)

		resp, err := s.client.SendContractTransaction(ctx, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
	
	s.T().Run("should fail with 422 if chain does not exist", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(404)
		txRequest := testutils.FakeSendTransactionRequest(chain.Name)

		resp, err := s.client.SendContractTransaction(ctx, txRequest)

		assert.Nil(t, resp)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Transactions() {
	ctx := context.Background()
	chain := testutils2.FakeChain()
	chainModel := &models.Chain{
		Name: chain.Name,
		UUID: chain.UUID,
		TenantID: chain.TenantID,
	}

	s.T().Run("should send a transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txRequest := testutils.FakeSendTransactionRequest(chain.Name)
	
		txResponse, err := s.client.SendContractTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, txRequest.IdempotencyKey, txResponse.IdempotencyKey)
		assert.Equal(t, txRequest.ChainName, chain.Name)
		assert.NotEmpty(t, txResponse.Schedule.UUID)
	
		scheduleResponse, err := s.client.GetSchedule(ctx, txResponse.Schedule.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, scheduleResponse.Jobs[0].UUID)
		assert.Equal(t, scheduleResponse.Jobs[0].ChainUUID, chain.UUID)
		assert.Equal(t, types.StatusStarted, scheduleResponse.Jobs[0].Status)
		assert.Equal(t, txRequest.Params.From, scheduleResponse.Jobs[0].Transaction.From)
		assert.Equal(t, txRequest.Params.To, scheduleResponse.Jobs[0].Transaction.To)
		assert.Equal(t, types.EthereumTransaction, scheduleResponse.Jobs[0].Type)

		evlp, err := s.env.consumer.WaitForEnvelope(scheduleResponse.Jobs[0].UUID, 
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, scheduleResponse.Jobs[0].UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[types.EthereumTransaction].String(), evlp.GetJobTypeString())
	})
	
	s.T().Run("should send a tessera transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txRequest := testutils.FakeSendTesseraRequest(chain.Name)
	
		txResponse, err := s.client.SendContractTransaction(ctx, txRequest)
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
	
		evlp, err := s.env.consumer.WaitForEnvelope(scheduleResponse.Jobs[0].UUID, 
			s.env.kafkaTopicConfig.Crafter,waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, scheduleResponse.Jobs[0].UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[types.TesseraPrivateTransaction].String(), evlp.GetJobTypeString())
	})
	
	s.T().Run("should send a orion transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txRequest := testutils.FakeSendOrionRequest(chain.Name)
	
		txResponse, err := s.client.SendContractTransaction(ctx, txRequest)
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
	
		evlp, err := s.env.consumer.WaitForEnvelope(scheduleResponse.Jobs[0].UUID, 
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, scheduleResponse.Jobs[0].UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[types.OrionEEATransaction].String(), evlp.GetJobTypeString())
	})
	
	s.T().Run("should send a deploy contract successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txRequest := testutils.FakeDeployContractRequest(chain.Name)
		txRequest.Params.Args = []string{"123"} // FakeContract arguments
	
		s.env.contractRegistryResponseFaker.GetContract = func() (*proto.GetContractResponse, error) {
			return &proto.GetContractResponse{
				Contract: testutils2.FakeContract(),
			}, nil
		}
		txResponse, err := s.client.SendDeployTransaction(ctx, txRequest)
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
	
		evlp, err := s.env.consumer.WaitForEnvelope(scheduleResponse.Jobs[0].UUID, 
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, scheduleResponse.Jobs[0].UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[types.EthereumTransaction].String(), evlp.GetJobTypeString())
	})
	
	s.T().Run("should send a raw transaction successfully to the transaction sender topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txRequest := testutils.FakeSendRawTransactionRequest(chain.Name)
	
		txResponse, err := s.client.SendRawTransaction(ctx, txRequest)
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
	
		evlp, err := s.env.consumer.WaitForEnvelope(scheduleResponse.Jobs[0].UUID, 
			s.env.kafkaTopicConfig.Sender, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, scheduleResponse.Jobs[0].UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[types.EthereumRawTransaction].String(), evlp.GetJobTypeString())
	})
	
	s.T().Run("should send a transfer transaction successfully to the transaction sender topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txRequest := testutils.FakeSendTransferTransactionRequest(chain.Name)

		txResponse, err := s.client.SendTransferTransaction(ctx, txRequest)
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
		assert.Equal(t, txRequest.Params.Value, scheduleResponse.Jobs[0].Transaction.Value)
		assert.Equal(t, txRequest.Params.To, scheduleResponse.Jobs[0].Transaction.To)
		assert.Equal(t, txRequest.Params.From, scheduleResponse.Jobs[0].Transaction.From)
		assert.Equal(t, types.EthereumTransaction, scheduleResponse.Jobs[0].Type)

		evlp, err := s.env.consumer.WaitForEnvelope(scheduleResponse.Jobs[0].UUID, 
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, scheduleResponse.Jobs[0].UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[types.EthereumTransaction].String(), evlp.GetJobTypeString())
	})
	
	s.T().Run("should succeed if payloads are the same and generate new schedule", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest(chain.Name)
	
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txResponse0, err := s.client.SendContractTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)
		txResponse1, err := s.client.SendContractTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		assert.NotEqual(t, txResponse0.Schedule.UUID, txResponse1.Schedule.UUID)
	})
}
