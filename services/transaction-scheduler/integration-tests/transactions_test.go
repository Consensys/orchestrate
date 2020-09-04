// +build integration

package integrationtests

import (
	"context"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/controllers"
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
	conf := client.NewConfig(s.baseURL, nil)
	s.client = client.NewHTTPClient(http.NewClient(), conf)
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Validation() {
	ctx := context.Background()
	chain := testutils.FakeChain()

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest()
		txRequest.ChainName = ""

		_, err := s.client.SendContractTransaction(ctx, txRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail if idempotency key is identical but different params", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest()
		rctx := context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			controllers.IdempotencyKeyHeader: utils.RandomString(16),
		})

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		_, err := s.client.SendContractTransaction(rctx, txRequest)
		assert.NoError(t, err)

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		txRequest.Params.MethodSignature = "differentMethodSignature()"
		_, err = s.client.SendContractTransaction(rctx, txRequest)
		assert.True(t, errors.IsConstraintViolatedError(err))
	})

	s.T().Run("should fail with 422 if chains cannot be fetched", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(404)
		txRequest := testutils.FakeSendTransactionRequest()

		_, err := s.client.SendContractTransaction(ctx, txRequest)

		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with 422 if chainUUID does not exist", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(404)
		txRequest := testutils.FakeSendTransactionRequest()

		_, err := s.client.SendContractTransaction(ctx, txRequest)

		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with 422 if from account is set and oneTimeKeyEnabled", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200)
		txRequest := testutils.FakeSendTransactionRequest()
		txRequest.Params.OneTimeKey = true

		_, err := s.client.SendContractTransaction(ctx, txRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Transactions() {
	ctx := context.Background()
	chain := testutils.FakeChain()

	s.T().Run("should send a transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		txRequest := testutils.FakeSendTransactionRequest()
		txRequest.Params.From = ""
		txRequest.Params.OneTimeKey = true
		IdempotencyKey := utils.RandomString(16)
		rctx := context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			controllers.IdempotencyKeyHeader: IdempotencyKey,
		})
		txResponse, err := s.client.SendContractTransaction(rctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, txResponse.UUID)
		assert.NotEmpty(t, txResponse.IdempotencyKey)

		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		job := txResponseGET.Schedule.Jobs[0]

		assert.NotEmpty(t, txResponseGET.Schedule.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, job.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.From, job.Transaction.From)
		assert.Equal(t, txRequest.Params.To, job.Transaction.To)
		assert.Equal(t, utils.EthereumTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.UUID, s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.UUID, evlp.GetID())
		assert.True(t, evlp.IsOneTimeKeySignature())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumTransaction].String(), evlp.GetJobTypeString())
		assert.Equal(t, evlp.GetChainIDString(), chain.ChainID)
		assert.Equal(t, evlp.PartitionKey(), "")
	})

	s.T().Run("should send a tessera transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
		txRequest := testutils.FakeSendTesseraRequest()
		IdempotencyKey := utils.RandomString(16)
		rctx := context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			controllers.IdempotencyKeyHeader: IdempotencyKey,
		})
		txResponse, err := s.client.SendContractTransaction(rctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, txResponse.UUID)
		assert.NotEmpty(t, txResponse.IdempotencyKey)

		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		
		assert.NotEmpty(t, txResponseGET.Schedule.UUID)
		assert.Len(t, txResponseGET.Schedule.Jobs, 2)
		
		privTxJob := txResponseGET.Schedule.Jobs[0]
		assert.NotEmpty(t, privTxJob.UUID)
		assert.Equal(t, privTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, privTxJob.Status)
		assert.Equal(t, txRequest.Params.From, privTxJob.Transaction.From)
		assert.Equal(t, txRequest.Params.To, privTxJob.Transaction.To)
		assert.Equal(t, utils.TesseraPrivateTransaction, privTxJob.Type)
		
		markingTxJob := txResponseGET.Schedule.Jobs[1]
		assert.NotEmpty(t, markingTxJob.UUID)
		assert.Equal(t, markingTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusCreated, markingTxJob.Status)
		assert.Equal(t, utils.TesseraMarkingTransaction, markingTxJob.Type)

		privTxEvlp, err := s.env.consumer.WaitForEnvelope(privTxJob.UUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		
		assert.Equal(t, privTxJob.UUID, privTxEvlp.GetID())
		assert.False(t, privTxEvlp.IsOneTimeKeySignature())
		assert.Equal(t, tx.JobTypeMap[utils.TesseraPrivateTransaction].String(), privTxEvlp.GetJobTypeString())
	})

	s.T().Run("should send an orion transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
		txRequest := testutils.FakeSendOrionRequest()
	
		txResponse, err := s.client.SendContractTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, txResponse.UUID)
		assert.NotEmpty(t, txResponse.IdempotencyKey)
	
		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		assert.NotEmpty(t, txResponseGET.Schedule.UUID)
		assert.Len(t, txResponseGET.Schedule.Jobs, 2)

		privTxJob := txResponseGET.Schedule.Jobs[0]
		assert.NotEmpty(t, privTxJob.UUID)
		assert.Equal(t, privTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, privTxJob.Status)
		assert.Equal(t, txRequest.Params.From, privTxJob.Transaction.From)
		assert.Equal(t, txRequest.Params.To, privTxJob.Transaction.To)
		assert.Equal(t, utils.OrionEEATransaction, privTxJob.Type)
		
		markingTxJob := txResponseGET.Schedule.Jobs[1]
		assert.NotEmpty(t, markingTxJob.UUID)
		assert.Equal(t, markingTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusCreated, markingTxJob.Status)
		assert.Equal(t, utils.OrionMarkingTransaction, markingTxJob.Type)
	
		privTxEvlp, err := s.env.consumer.WaitForEnvelope(privTxJob.UUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		
		assert.Equal(t, privTxJob.UUID, privTxEvlp.GetID())
		assert.Equal(t, tx.JobTypeMap[utils.OrionEEATransaction].String(), privTxEvlp.GetJobTypeString())
	})

	s.T().Run("should send a deploy contract successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		txRequest := testutils.FakeDeployContractRequest()
		txRequest.Params.Args = testutils.ParseIArray(123) // FakeContract arguments

		s.env.contractRegistryResponseFaker.GetContract = func() (*proto.GetContractResponse, error) {
			return &proto.GetContractResponse{
				Contract: testutils.FakeContract(),
			}, nil
		}
		txResponse, err := s.client.SendDeployTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, txResponse.UUID)
		assert.NotEmpty(t, txResponse.IdempotencyKey)

		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		job := txResponseGET.Schedule.Jobs[0]

		assert.NotEmpty(t, txResponseGET.Schedule.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, job.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.From, job.Transaction.From)
		assert.Equal(t, utils.EthereumTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.UUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumTransaction].String(), evlp.GetJobTypeString())
	})

	s.T().Run("should send a raw transaction successfully to the transaction sender topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		txRequest := testutils.FakeSendRawTransactionRequest()

		IdempotencyKey := utils.RandomString(16)
		rctx := context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			controllers.IdempotencyKeyHeader: IdempotencyKey,
		})
		txResponse, err := s.client.SendRawTransaction(rctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.NotEmpty(t, txResponse.UUID)
		assert.NotEmpty(t, txResponse.IdempotencyKey)

		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		job := txResponseGET.Schedule.Jobs[0]

		assert.NotEmpty(t, txResponseGET.Schedule.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.Raw, job.Transaction.Raw)
		assert.Equal(t, utils.EthereumRawTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.UUID,
			s.env.kafkaTopicConfig.Sender, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumRawTransaction].String(), evlp.GetJobTypeString())
	})

	s.T().Run("should send a transfer transaction successfully to the transaction sender topic", func(t *testing.T) {
		defer gock.Off()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		txRequest := testutils.FakeSendTransferTransactionRequest()

		txResponse, err := s.client.SendTransferTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Len(t, txResponse.IdempotencyKey, 16)
		assert.NotEmpty(t, txResponse.UUID)

		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		job := txResponseGET.Schedule.Jobs[0]

		assert.NotEmpty(t, txResponseGET.Schedule.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.Value, job.Transaction.Value)
		assert.Equal(t, txRequest.Params.To, job.Transaction.To)
		assert.Equal(t, txRequest.Params.From, job.Transaction.From)
		assert.Equal(t, utils.EthereumTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.UUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.UUID, evlp.GetID())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumTransaction].String(), evlp.GetJobTypeString())
	})

	s.T().Run("should succeed if payloads and idempotency key are the same and return same schedule", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest()
		rctx := context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			controllers.IdempotencyKeyHeader: utils.RandomString(16),
		})

		// Kill Kafka on first call so data is added in DB and status is CREATED but does not get updated to STARTED
		err := s.env.client.Stop(rctx, kafkaContainerID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		_, err = s.client.SendContractTransaction(rctx, txRequest)
		assert.Error(t, err)

		s.restartKafka(rctx)

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		txResponse, err := s.client.SendContractTransaction(rctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		job := txResponse.Schedule.Jobs[0]
		assert.Equal(t, utils.StatusStarted, job.Status)
	})
}

func (s *txSchedulerTransactionTestSuite) restartKafka(ctx context.Context) {
	// Bring Kafka Up and resend so the transaction is sent successfully
	err := s.env.client.Start(ctx, kafkaContainerID)
	if err != nil {
		s.env.logger.WithError(err).Error("could not restart kafka")
		return
	}
	err = s.env.client.WaitTillIsReady(ctx, kafkaContainerID, 20*time.Second)
	if err != nil {
		s.env.logger.WithError(err).Error("could not start transaction-scheduler")
		return
	}
}
