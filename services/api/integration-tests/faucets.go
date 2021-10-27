// +build integration

package integrationtests

import (
	"testing"
	"time"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"

	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type faucetsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *faucetsTestSuite) TestRegister() {
	ctx := s.env.ctx

	s.T().Run("should register faucet successfully", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()

		resp, err := s.client.RegisterFaucet(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, req.CreditorAccount, resp.CreditorAccount)
		assert.Equal(t, req.ChainRule, resp.ChainRule)
		assert.Equal(t, req.MaxBalance, resp.MaxBalance)
		assert.Equal(t, req.Amount, resp.Amount)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Cooldown, resp.Cooldown)
		assert.NotEmpty(t, resp.UUID)
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.Equal(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteFaucet(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should fail to register faucet with same name and tenant", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()

		resp, err := s.client.RegisterFaucet(ctx, req)
		require.NoError(t, err)

		_, err = s.client.RegisterFaucet(ctx, req)
		assert.True(t, errors.IsAlreadyExistsError(err))

		err = s.client.DeleteFaucet(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should fail to register faucet if postgres is down", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		_, err = s.client.RegisterFaucet(ctx, req)
		assert.Error(t, err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(t, err)
	})
}

func (s *faucetsTestSuite) TestSearch() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterFaucetRequest()
	faucet, err := s.client.RegisterFaucet(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should search faucet by name successfully", func(t *testing.T) {
		resp, err := s.client.SearchFaucets(ctx, &entities.FaucetFilters{
			Names: []string{faucet.Name},
		})
		require.NoError(t, err)

		assert.Len(t, resp, 1)
		assert.Equal(t, faucet.UUID, resp[0].UUID)
	})

	s.T().Run("should search faucet by chain_rule successfully", func(t *testing.T) {
		resp, err := s.client.SearchFaucets(ctx, &entities.FaucetFilters{
			ChainRule: faucet.ChainRule,
		})
		require.NoError(t, err)

		assert.Len(t, resp, 1)
		assert.Equal(t, faucet.UUID, resp[0].UUID)
	})

	err = s.client.DeleteFaucet(ctx, faucet.UUID)
	require.NoError(s.T(), err)
}

func (s *faucetsTestSuite) TestGetOne() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterFaucetRequest()
	faucet, err := s.client.RegisterFaucet(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should get faucet successfully", func(t *testing.T) {
		resp, err := s.client.GetFaucet(ctx, faucet.UUID)
		require.NoError(t, err)
		assert.Equal(t, faucet.UUID, resp.UUID)
	})

	err = s.client.DeleteFaucet(ctx, faucet.UUID)
	require.NoError(s.T(), err)
}

func (s *faucetsTestSuite) TestUpdate() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterFaucetRequest()
	faucet, err := s.client.RegisterFaucet(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should update faucet successfully", func(t *testing.T) {
		req := testutils.FakeUpdateFaucetRequest()

		resp, err := s.client.UpdateFaucet(ctx, faucet.UUID, req)
		require.NoError(t, err)

		assert.Equal(t, req.CreditorAccount, resp.CreditorAccount)
		assert.Equal(t, req.ChainRule, resp.ChainRule)
		assert.Equal(t, req.MaxBalance, resp.MaxBalance)
		assert.Equal(t, req.Amount, resp.Amount)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Cooldown, resp.Cooldown)
		assert.NotEmpty(t, resp.UUID)
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.NotEqual(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteFaucet(ctx, resp.UUID)
		assert.NoError(t, err)
	})
}

func (s *faucetsTestSuite) TestSuccess_TxsWithFaucet() {
	ctx := s.env.ctx

	chainWithFaucet, err := s.client.RegisterChain(s.env.ctx, &api.RegisterChainRequest{
		Name: "ganache-with-faucet",
		URLs: []string{s.env.blockchainNodeURL},
		Listener: api.RegisterListenerRequest{
			FromBlock:         "latest",
			ExternalTxEnabled: false,
		},
	})
	require.NoError(s.T(), err)

	accountFaucetAlias := "MyFaucetCreditor"
	req := testutils.FakeImportAccountRequest()
	req.Alias = accountFaucetAlias
	// Ganache imported account with 1000ETH
	req.PrivateKey = hexutil.MustDecode("0x56202652fdffd802b7252a456dbd8f3ecc0352bbde76c23b40afe8aebd714e2e")
	accResp, err := s.client.ImportAccount(s.env.ctx, req)
	require.NoError(s.T(), err)

	faucetRequest := testutils.FakeRegisterFaucetRequest()
	faucetRequest.Name = "faucet-integration-tests"
	faucetRequest.ChainRule = chainWithFaucet.UUID
	faucetRequest.CreditorAccount = accResp.Address
	faucetRequest.Cooldown = "0s"
	faucet, err := s.client.RegisterFaucet(s.env.ctx, faucetRequest)
	require.NoError(s.T(), err)
	
	defer func() {
		err = s.client.DeleteChain(ctx, chainWithFaucet.UUID)
		assert.NoError(s.T(), err)
		err = s.client.DeleteFaucet(ctx, faucet.UUID)
		assert.NoError(s.T(), err)
	}()

	s.T().Run("should send a transaction with an additional faucet job", func(t *testing.T) {
		defer gock.Off()
		// Transfer tx
		txRequest := testutils.FakeSendTransferTransactionRequest()
		txRequest.ChainName = chainWithFaucet.Name
		txResponse, err := s.client.SendTransferTransaction(ctx, txRequest)
		require.NoError(t, err)
		assert.NotEmpty(t, txResponse.UUID)
	
		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		require.NoError(t, err)
		require.Len(t, txResponseGET.Jobs, 2)
	
		faucetJob := txResponseGET.Jobs[1]
		txJob := txResponseGET.Jobs[0]
		assert.Equal(t, faucetJob.ChainUUID, faucet.ChainRule)
		assert.Equal(t, entities.StatusStarted, faucetJob.Status)
		assert.Equal(t, entities.EthereumTransaction, faucetJob.Type)
		assert.Equal(t, faucetJob.Transaction.To, txJob.Transaction.From)
		assert.Equal(t, faucetJob.Transaction.Value, faucet.Amount)
	
		assert.NotEmpty(t, txResponseGET.UUID)
		assert.NotEmpty(t, txJob.UUID)
		assert.Equal(t, txJob.ChainUUID, faucet.ChainRule)
		assert.Equal(t, entities.StatusStarted, txJob.Status)
		assert.Equal(t, txRequest.Params.From, txJob.Transaction.From)
		assert.Equal(t, txRequest.Params.To, txJob.Transaction.To)
		assert.Equal(t, entities.EthereumTransaction, txJob.Type)
	
		fctEvlp, err := s.env.consumer.WaitForEnvelope(faucetJob.ScheduleUUID, s.env.kafkaTopicConfig.Sender, waitForEnvelopeTimeOut)
		require.NoError(t, err)
		assert.Equal(t, faucetJob.ScheduleUUID, fctEvlp.GetID())
		assert.Equal(t, faucetJob.UUID, fctEvlp.GetJobUUID())
	
		jobEvlp, err := s.env.consumer.WaitForEnvelope(txJob.ScheduleUUID, s.env.kafkaTopicConfig.Sender, waitForEnvelopeTimeOut)
		require.NoError(t, err)
		assert.Equal(t, txJob.ScheduleUUID, jobEvlp.GetID())
		assert.Equal(t, txJob.UUID, jobEvlp.GetJobUUID())
	})

	s.T().Run("should send a raw transaction with an additional faucet job", func(t *testing.T) {
		defer gock.Off()
		// Raw tx
		txRequest := testutils.FakeSendRawTransactionRequest()
		txRequest.ChainName = chainWithFaucet.Name
		txResponse, err := s.client.SendRawTransaction(ctx, txRequest)
		require.NoError(t, err)
		assert.NotEmpty(t, txResponse.UUID)

		txResponseGET, err := s.client.GetTxRequest(ctx, txResponse.UUID)
		require.NoError(t, err)
		require.Len(t, txResponseGET.Jobs, 2)

		faucetJob := txResponseGET.Jobs[1]
		txJob := txResponseGET.Jobs[0]
		assert.Equal(t, faucetJob.ChainUUID, faucet.ChainRule)
		assert.Equal(t, entities.StatusStarted, faucetJob.Status)
		assert.Equal(t, entities.EthereumTransaction, faucetJob.Type)
		assert.Equal(t, faucetJob.Transaction.To, txJob.Transaction.From)
		assert.Equal(t, faucetJob.Transaction.Value, faucet.Amount)

		assert.NotEmpty(t, txResponseGET.UUID)
		assert.NotEmpty(t, txJob.UUID)
		assert.Equal(t, txJob.ChainUUID, faucet.ChainRule)
		assert.Equal(t, entities.StatusStarted, txJob.Status)
		assert.Equal(t, "0x4c7aF4B315644848f400b7344A8e73Cf227812b4", txJob.Transaction.From)
		assert.Equal(t, entities.EthereumRawTransaction, txJob.Type)

		fctEvlp, err := s.env.consumer.WaitForEnvelope(faucetJob.ScheduleUUID, s.env.kafkaTopicConfig.Sender, waitForEnvelopeTimeOut)
		require.NoError(t, err)
		assert.Equal(t, faucetJob.ScheduleUUID, fctEvlp.GetID())
		assert.Equal(t, faucetJob.UUID, fctEvlp.GetJobUUID())

		jobEvlp, err := s.env.consumer.WaitForEnvelope(txJob.ScheduleUUID, s.env.kafkaTopicConfig.Sender, waitForEnvelopeTimeOut)
		require.NoError(t, err)
		assert.Equal(t, txJob.ScheduleUUID, jobEvlp.GetID())
		assert.Equal(t, txJob.UUID, jobEvlp.GetJobUUID())
	})
}
