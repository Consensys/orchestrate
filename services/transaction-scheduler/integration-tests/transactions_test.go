// +build integration

package integrationtests

import (
	"context"
	"fmt"
	http2 "net/http"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client-utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
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
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_Validation() {
	ctx := s.env.ctx
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
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()
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
		txRequest := testutils.FakeSendTransactionRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(404)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()

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
	ctx := s.env.ctx
	chain := testutils.FakeChain()
	faucet := testutils.FakeFaucet()

	s.T().Run("should send a transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransactionRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, ethcommon.HexToAddress("0x").Hex())).
			Reply(404).Done()
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

		job := txResponseGET.Jobs[0]

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, job.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.From, job.Transaction.From)
		assert.Equal(t, txRequest.Params.To, job.Transaction.To)
		assert.Equal(t, utils.EthereumTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.ScheduleUUID, s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.ScheduleUUID, evlp.GetID())
		assert.Equal(t, job.UUID, evlp.GetJobUUID())
		assert.True(t, evlp.IsOneTimeKeySignature())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumTransaction].String(), evlp.GetJobTypeString())
		assert.Equal(t, evlp.GetChainIDString(), chain.ChainID)
		assert.Equal(t, evlp.PartitionKey(), "")
	})

	s.T().Run("should send a transaction, with an additional faucet job, successfully to the transaction crafter topic in parallel", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransferTransactionRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(200).JSON(faucet)
		IdempotencyKey := utils.RandomString(16)
		rctx := context.WithValue(ctx, clientutils.RequestHeaderKey, map[string]string{
			controllers.IdempotencyKeyHeader: IdempotencyKey,
		})
		txResponse, err := s.client.SendTransferTransaction(rctx, txRequest)
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

		faucetJob := txResponseGET.Jobs[1]
		txJob := txResponseGET.Jobs[0]
		assert.Equal(t, faucetJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, faucetJob.Status)
		assert.Equal(t, utils.EthereumTransaction, faucetJob.Type)
		assert.Equal(t, faucetJob.Transaction.To, txJob.Transaction.From)
		assert.Equal(t, faucetJob.Transaction.Value, faucet.Amount.String())

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.NotEmpty(t, txJob.UUID)
		assert.Equal(t, txJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, txJob.Status)
		assert.Equal(t, txRequest.Params.From, txJob.Transaction.From)
		assert.Equal(t, txRequest.Params.To, txJob.Transaction.To)
		assert.Equal(t, utils.EthereumTransaction, txJob.Type)

		fctEvlp, err := s.env.consumer.WaitForEnvelope(faucetJob.ScheduleUUID, s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, faucetJob.ScheduleUUID, fctEvlp.GetID())
		assert.Equal(t, faucetJob.UUID, fctEvlp.GetJobUUID())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumTransaction].String(), fctEvlp.GetJobTypeString())
		assert.Equal(t, fctEvlp.GetChainIDString(), chain.ChainID)

		jobEvlp, err := s.env.consumer.WaitForEnvelope(txJob.ScheduleUUID, s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, txJob.ScheduleUUID, jobEvlp.GetID())
		assert.Equal(t, txJob.UUID, jobEvlp.GetJobUUID())
	})

	s.T().Run("should send a tessera transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTesseraRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()
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

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.Len(t, txResponseGET.Jobs, 2)

		privTxJob := txResponseGET.Jobs[0]
		assert.NotEmpty(t, privTxJob.UUID)
		assert.Equal(t, privTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, privTxJob.Status)
		assert.Equal(t, txRequest.Params.From, privTxJob.Transaction.From)
		assert.Equal(t, txRequest.Params.To, privTxJob.Transaction.To)
		assert.Equal(t, utils.TesseraPrivateTransaction, privTxJob.Type)

		markingTxJob := txResponseGET.Jobs[1]
		assert.NotEmpty(t, markingTxJob.UUID)
		assert.Equal(t, markingTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusCreated, markingTxJob.Status)
		assert.Equal(t, utils.TesseraMarkingTransaction, markingTxJob.Type)

		privTxEvlp, err := s.env.consumer.WaitForEnvelope(privTxJob.ScheduleUUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, privTxJob.ScheduleUUID, privTxEvlp.GetID())
		assert.Equal(t, privTxJob.UUID, privTxEvlp.GetJobUUID())
		assert.False(t, privTxEvlp.IsOneTimeKeySignature())
		assert.Equal(t, tx.JobTypeMap[utils.TesseraPrivateTransaction].String(), privTxEvlp.GetJobTypeString())
	})

	s.T().Run("should send an orion transaction successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendOrionRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()

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

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.Len(t, txResponseGET.Jobs, 2)

		privTxJob := txResponseGET.Jobs[0]
		assert.NotEmpty(t, privTxJob.UUID)
		assert.Equal(t, privTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, privTxJob.Status)
		assert.Equal(t, txRequest.Params.From, privTxJob.Transaction.From)
		assert.Equal(t, txRequest.Params.To, privTxJob.Transaction.To)
		assert.Equal(t, utils.OrionEEATransaction, privTxJob.Type)

		markingTxJob := txResponseGET.Jobs[1]
		assert.NotEmpty(t, markingTxJob.UUID)
		assert.Equal(t, markingTxJob.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusCreated, markingTxJob.Status)
		assert.Equal(t, utils.OrionMarkingTransaction, markingTxJob.Type)

		privTxEvlp, err := s.env.consumer.WaitForEnvelope(privTxJob.ScheduleUUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, privTxJob.ScheduleUUID, privTxEvlp.GetID())
		assert.Equal(t, privTxJob.UUID, privTxEvlp.GetJobUUID())
		assert.Equal(t, tx.JobTypeMap[utils.OrionEEATransaction].String(), privTxEvlp.GetJobTypeString())
	})

	s.T().Run("should send a deploy contract successfully to the transaction crafter topic", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeDeployContractRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()
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

		job := txResponseGET.Jobs[0]

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, job.ChainUUID, chain.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.From, job.Transaction.From)
		assert.Equal(t, utils.EthereumTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.ScheduleUUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.ScheduleUUID, evlp.GetID())
		assert.Equal(t, job.UUID, evlp.GetJobUUID())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumTransaction].String(), evlp.GetJobTypeString())
	})

	s.T().Run("should send a raw transaction successfully to the transaction sender topic", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendRawTransactionRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)

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

		job := txResponseGET.Jobs[0]

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.Raw, job.Transaction.Raw)
		assert.Equal(t, utils.EthereumRawTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.ScheduleUUID,
			s.env.kafkaTopicConfig.Sender, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.ScheduleUUID, evlp.GetID())
		assert.Equal(t, job.UUID, evlp.GetJobUUID())
		assert.Equal(t, tx.JobTypeMap[utils.EthereumRawTransaction].String(), evlp.GetJobTypeString())
	})

	s.T().Run("should send a transfer transaction successfully to the transaction sender topic", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeSendTransferTransactionRequest()
		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()

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

		job := txResponseGET.Jobs[0]

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.NotEmpty(t, job.UUID)
		assert.Equal(t, utils.StatusStarted, job.Status)
		assert.Equal(t, txRequest.Params.Value, job.Transaction.Value)
		assert.Equal(t, txRequest.Params.To, job.Transaction.To)
		assert.Equal(t, txRequest.Params.From, job.Transaction.From)
		assert.Equal(t, utils.EthereumTransaction, job.Type)

		evlp, err := s.env.consumer.WaitForEnvelope(job.ScheduleUUID,
			s.env.kafkaTopicConfig.Crafter, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		assert.Equal(t, job.ScheduleUUID, evlp.GetID())
		assert.Equal(t, job.UUID, evlp.GetJobUUID())
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
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()
		_, err = s.client.SendContractTransaction(rctx, txRequest)
		assert.Error(t, err)

		err = s.env.client.StartServiceAndWait(rctx, kafkaContainerID, 10*time.Second)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		gock.New(ChainRegistryURL).
			URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", ChainRegistryURL, chain.UUID, txRequest.Params.From)).
			Reply(404).Done()
		txResponse, err := s.client.SendContractTransaction(rctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		job := txResponse.Jobs[0]
		assert.Equal(t, utils.StatusStarted, job.Status)
	})
}

func (s *txSchedulerTransactionTestSuite) TestTransactionScheduler_ZHealthCheck() {
	type healthRes struct {
		ChainRegistry    string `json:"chain-registry,omitempty"`
		ContractRegistry string `json:"contract-registry,omitempty"`
		Database         string `json:"Database,omitempty"`
		Kafka            string `json:"Kafka,omitempty"`
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	ctx := s.env.ctx
	s.T().Run("should retrieve positive health check over service dependencies", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(ChainRegistryURL).Get("/live").Reply(200)
		gock.New(ContractRegistryURL).Get("/live").Reply(200)
		defer gock.Off()

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), 200, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "OK", status.Database)
		assert.Equal(s.T(), "OK", status.ChainRegistry)
		assert.Equal(s.T(), "OK", status.ContractRegistry)
		assert.Equal(s.T(), "OK", status.Kafka)
	})

	s.T().Run("should retrieve a negative health check over chain-registry and contract-registry services ", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(ChainRegistryURL).Get("/live").Reply(500)
		gock.New(ContractRegistryURL).Get("/live").Reply(500)
		defer gock.Off()

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), 503, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "OK", status.Database)
		assert.Equal(s.T(), "OK", status.Kafka)
		assert.NotEqual(s.T(), "OK", status.ChainRegistry)
		assert.NotEqual(s.T(), "OK", status.ContractRegistry)
	})

	s.T().Run("should retrieve a negative health check over kafka service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(ChainRegistryURL).Get("/live").Reply(200)
		gock.New(ContractRegistryURL).Get("/live").Reply(200)
		defer gock.Off()

		// Kill Kafka on first call so data is added in DB and status is CREATED but does not get updated to STARTED
		err = s.env.client.Stop(ctx, kafkaContainerID)
		assert.NoError(t, err)

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		err = s.env.client.StartServiceAndWait(ctx, kafkaContainerID, 10*time.Second)
		assert.NoError(t, err)

		assert.Equal(s.T(), 503, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "OK", status.Database)
		assert.NotEqual(s.T(), "OK", status.Kafka)
		assert.Equal(s.T(), "OK", status.ChainRegistry)
		assert.Equal(s.T(), "OK", status.ContractRegistry)
	})

	s.T().Run("should retrieve a negative health check over postgres service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(ChainRegistryURL).Get("/live").Reply(200)
		gock.New(ContractRegistryURL).Get("/live").Reply(200)
		defer gock.Off()

		// Kill Kafka on first call so data is added in DB and status is CREATED but does not get updated to STARTED
		err = s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), 503, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.NotEqual(s.T(), "OK", status.Database)
		assert.Equal(s.T(), "OK", status.Kafka)
		assert.Equal(s.T(), "OK", status.ChainRegistry)
		assert.Equal(s.T(), "OK", status.ContractRegistry)
	})
}
