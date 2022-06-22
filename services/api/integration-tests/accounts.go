// +build integration

package integrationtests

import (
	"github.com/consensys/orchestrate/pkg/ethereum/account"
	utilstypes "github.com/consensys/quorum-key-manager/src/utils/api/types"
	"testing"
	"time"

	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/types/api"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/traefik/traefik/v2/pkg/log"
)

type accountsTestSuite struct {
	suite.Suite
	client            client.OrchestrateClient
	env               *IntegrationEnvironment
	defaultQKMStoreID string
}

func (s *accountsTestSuite) TestCreate() {
	ctx := s.env.ctx
	chain := testutils.FakeChain()
	chain.URLs = []string{s.env.blockchainNodeURL}

	s.T().Run("should create account successfully by querying key-manager API", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
		require.NoError(s.T(), err)

		resp, err := s.client.GetAccount(ctx, ethAccRes.Address)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), resp.Address, ethAccRes.Address)
		assert.Equal(s.T(), resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(s.T(), resp.Alias, txRequest.Alias)
		assert.Equal(s.T(), resp.StoreID, s.defaultQKMStoreID)
		assert.Equal(s.T(), resp.TenantID, "_")
	})

	s.T().Run("should fail to create account if QKM storeID does not exist", func(t *testing.T) {
		qkmStoreID := "my-personal-storeID"
		txRequest := testutils.FakeCreateAccountRequest()
		txRequest.StoreID = qkmStoreID

		_, err := s.client.CreateAccount(ctx, txRequest)
		require.Error(s.T(), err)
		// QKM StoreID does not exist
		require.True(s.T(), errors.IsDependencyFailureError(err))
	})

	s.T().Run("should fail to create account with same alias", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		_, err := s.client.CreateAccount(ctx, txRequest)
		require.NoError(s.T(), err)

		_, err = s.client.CreateAccount(ctx, txRequest)
		assert.Error(s.T(), err)
		log.WithoutContext().Errorf("%v", err)
		assert.True(s.T(), errors.IsAlreadyExistsError(err))
	})

	s.T().Run("should create account successfully and trigger funding transaction", func(t *testing.T) {
		chainWithFaucet, err := s.client.RegisterChain(s.env.ctx, &api.RegisterChainRequest{
			Name: "ganache-with-faucet-accounts",
			URLs: []string{s.env.blockchainNodeURL},
			Listener: api.RegisterListenerRequest{
				FromBlock:         "latest",
			},
		})
		require.NoError(s.T(), err)

		acc, err := account.NewAccount()
		require.NoError(s.T(), err)
		accountFaucet := testutils.FakeAccount()
		accountFaucet.Alias = "MyFaucetCreditor-accounts_" + utils.RandString(5)
		accountFaucet.Address = acc.Address

		req := testutils.FakeImportAccountRequest()
		req.PrivateKey = acc.Priv()
		req.Alias = accountFaucet.Alias
		ethAccRes, err := s.client.ImportAccount(s.env.ctx, req)
		require.NoError(s.T(), err)

		faucetRequest := testutils.FakeRegisterFaucetRequest()
		faucetRequest.Name = "faucet-integration-tests"
		faucetRequest.ChainRule = chainWithFaucet.UUID
		faucetRequest.CreditorAccount = accountFaucet.Address
		faucet, err := s.client.RegisterFaucet(s.env.ctx, faucetRequest)
		require.NoError(s.T(), err)

		accountRequest := testutils.FakeCreateAccountRequest()
		accountRequest.Chain = chainWithFaucet.Name

		require.NoError(s.T(), err)

		assert.Equal(s.T(), ethAccRes.TenantID, "_")

		err = s.client.DeleteChain(ctx, chainWithFaucet.UUID)
		assert.NoError(s.T(), err)
		err = s.client.DeleteFaucet(ctx, faucet.UUID)
		assert.NoError(s.T(), err)
	})

	s.T().Run("should fail to create account if postgres is down", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(s.T(), err)

		_, err = s.client.CreateAccount(ctx, txRequest)
		assert.Error(s.T(), err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(s.T(), err)
	})
}

func (s *accountsTestSuite) TestImport() {
	ctx := s.env.ctx

	s.T().Run("should import account successfully by querying key-manager API", func(t *testing.T) {
		acc, err := account.NewAccount()
		require.NoError(s.T(), err)
		txRequest := testutils.FakeImportAccountRequest()
		txRequest.PrivateKey = acc.Priv()

		resp, err := s.client.ImportAccount(ctx, txRequest)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), resp.Alias, txRequest.Alias)
		assert.Equal(s.T(), resp.TenantID, "_")
	})

	s.T().Run("should fail to import same account twice", func(t *testing.T) {
		acc, err := account.NewAccount()
		require.NoError(s.T(), err)

		txRequest := testutils.FakeImportAccountRequest()
		txRequest.PrivateKey = acc.Priv()

		_, err = s.client.ImportAccount(ctx, txRequest)
		require.NoError(s.T(), err)

		_, err = s.client.ImportAccount(ctx, txRequest)
		require.Error(s.T(), err)
		assert.True(s.T(), errors.IsAlreadyExistsError(err))
	})
}

func (s *accountsTestSuite) TestSearch() {
	ctx := s.env.ctx

	s.T().Run("should create account and search for it by alias successfully", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
		require.NoError(s.T(), err)
		resp, err := s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: []string{txRequest.Alias},
		})
		require.NoError(s.T(), err)

		assert.Len(s.T(), resp.Accounts, 1)
		assert.Equal(s.T(), resp.Accounts[0].Address, ethAccRes.Address)
		assert.Equal(s.T(), resp.Accounts[0].PublicKey, ethAccRes.PublicKey)
		assert.Equal(s.T(), resp.Accounts[0].Alias, txRequest.Alias)
		assert.Equal(s.T(), resp.Accounts[0].TenantID, "_")
	})

	s.T().Run("should create account and search with limit and page successfully", func(t *testing.T) {
		txRequest0 := testutils.FakeCreateAccountRequest()

		ethAccRes0, err := s.client.CreateAccount(ctx, txRequest0)
		require.NoError(s.T(), err)
		txRequest1 := testutils.FakeCreateAccountRequest()

		ethAccRes1, err := s.client.CreateAccount(ctx, txRequest1)
		resp, err := s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: []string{txRequest0.Alias, txRequest1.Alias},
			Pagination: entities.PaginationFilters{Limit: 1},
		})
		require.NoError(s.T(), err)

		assert.Len(s.T(), resp.Accounts, 1)
		assert.Equal(s.T(), resp.Accounts[0].Address, ethAccRes0.Address)
		assert.Equal(s.T(), resp.Accounts[0].PublicKey, ethAccRes0.PublicKey)
		assert.Equal(s.T(), resp.Accounts[0].Alias, txRequest0.Alias)
		assert.Equal(s.T(), resp.Accounts[0].TenantID, "_")

		resp, err = s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: []string{txRequest0.Alias, txRequest1.Alias},
			Pagination: entities.PaginationFilters{Limit: 2},
		})
		require.NoError(s.T(), err)

		assert.Len(s.T(), resp.Accounts, 2)
		assert.Equal(s.T(), resp.Accounts[1].Address, ethAccRes1.Address)
		assert.Equal(s.T(), resp.Accounts[1].PublicKey, ethAccRes1.PublicKey)
		assert.Equal(s.T(), resp.Accounts[1].Alias, txRequest1.Alias)
		assert.Equal(s.T(), resp.Accounts[1].TenantID, "_")
	})

	s.T().Run("should create accounts and search with limit and page successfully", func(t *testing.T) {

		txRequests := [25]*api.CreateAccountRequest{}
		ethAccounts := [25]*api.AccountResponse{}
		var aliases []string
		var err error

		for i,_ := range txRequests{
			txRequests[i]  = testutils.FakeCreateAccountRequest()
			ethAccounts[i], err = s.client.CreateAccount(ctx, txRequests[i])
			require.NoError(s.T(), err)
			aliases = append(aliases, txRequests[i].Alias)
		}

		resp, err := s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: aliases,
			Pagination: entities.PaginationFilters{Limit: 25},
		})
		require.NoError(s.T(), err)

		assert.Len(s.T(), resp.Accounts, 25)
		assert.Equal(s.T(), resp.Accounts[0].Address, ethAccounts[0].Address)
		assert.Equal(s.T(), resp.Accounts[0].PublicKey, ethAccounts[0].PublicKey)
		assert.Equal(s.T(), resp.Accounts[0].Alias, txRequests[0].Alias)
		assert.Equal(s.T(), resp.Accounts[0].TenantID, "_")

		resp, err = s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: []string{txRequests[0].Alias, txRequests[1].Alias},
			Pagination: entities.PaginationFilters{Limit: 2},
		})
		require.NoError(s.T(), err)

		assert.Len(s.T(), resp.Accounts, 2)
		assert.Equal(s.T(), resp.Accounts[1].Address, ethAccounts[1].Address)
		assert.Equal(s.T(), resp.Accounts[1].PublicKey, ethAccounts[1].PublicKey)
		assert.Equal(s.T(), resp.Accounts[1].Alias, txRequests[1].Alias)
		assert.Equal(s.T(), resp.Accounts[1].TenantID, "_")

		resp, err = s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: aliases,
			Pagination: entities.PaginationFilters{Limit: 5, Page: 1},
		})

		assert.Len(s.T(), resp.Accounts, 5)
		assert.Equal(s.T(), resp.Accounts[4].Address, ethAccounts[9].Address)
		assert.Equal(s.T(), resp.Accounts[4].PublicKey, ethAccounts[9].PublicKey)
		assert.Equal(s.T(), resp.Accounts[4].Alias, txRequests[9].Alias)
		assert.Equal(s.T(), resp.HasMore, true)
		assert.Equal(s.T(), resp.Accounts[4].TenantID, "_")

		resp, err = s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: aliases,
			Pagination: entities.PaginationFilters{Limit: 5, Page: 2},
		})

		assert.Len(s.T(), resp.Accounts, 5)
		assert.Equal(s.T(), resp.Accounts[4].Address, ethAccounts[14].Address)
		assert.Equal(s.T(), resp.Accounts[4].PublicKey, ethAccounts[14].PublicKey)
		assert.Equal(s.T(), resp.Accounts[4].Alias, txRequests[14].Alias)
		assert.Equal(s.T(), resp.HasMore, true)
		assert.Equal(s.T(), resp.Accounts[4].TenantID, "_")

		resp, err = s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: aliases,
			Pagination: entities.PaginationFilters{Limit: 5, Page: 4},
		})
		assert.Equal(s.T(), resp.HasMore, false)
	})
}

func (s *accountsTestSuite) TestUpdate() {
	ctx := s.env.ctx

	s.T().Run("should create account and update it successfully", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
		require.NoError(s.T(), err)

		txRequest2 := testutils.FakeUpdateAccountRequest()
		resp, err := s.client.UpdateAccount(ctx, ethAccRes.Address, txRequest2)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), resp.Alias, txRequest2.Alias)
		assert.Equal(s.T(), resp.Attributes, txRequest2.Attributes)
		assert.Equal(s.T(), resp.TenantID, "_")
	})
}

func (s *accountsTestSuite) TestSignMessageAndVerify() {
	ctx := s.env.ctx
	txRequest := testutils.FakeCreateAccountRequest()
	ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
	require.NoError(s.T(), err)

	message := hexutil.MustDecode("0xaeff")
	var signedPayload string

	s.T().Run("should sign message successfully", func(t *testing.T) {
		signedPayload, err = s.client.SignMessage(ctx, ethAccRes.Address, &qkmtypes.SignMessageRequest{
			Message: message,
		})
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedPayload)
	})

	s.T().Run("should verify signature successfully", func(t *testing.T) {
		verifyRequest := &utilstypes.VerifyRequest{
			Data:      message,
			Signature: hexutil.MustDecode(signedPayload),
			Address:   ethAccRes.Address,
		}
		err := s.client.VerifyMessageSignature(ctx, verifyRequest)
		assert.NoError(s.T(), err)
	})
}

func (s *accountsTestSuite) TestSignTypedData() {
	ctx := s.env.ctx

	txRequest := testutils.FakeCreateAccountRequest()
	ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
	require.NoError(s.T(), err)

	typedDataRequest := qkm.FakeSignTypedDataRequest()
	var signature string

	s.T().Run("should sign typed data successfully", func(t *testing.T) {
		signature, err = s.client.SignTypedData(ctx, ethAccRes.Address, &qkmtypes.SignTypedDataRequest{
			DomainSeparator: typedDataRequest.DomainSeparator,
			Types:           typedDataRequest.Types,
			Message:         typedDataRequest.Message,
			MessageType:     typedDataRequest.MessageType,
		})

		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signature)
	})

	s.T().Run("should verify typed data signature successfully", func(t *testing.T) {
		err := s.client.VerifyTypedDataSignature(ctx, &utilstypes.VerifyTypedDataRequest{
			TypedData: *typedDataRequest,
			Signature: hexutil.MustDecode(signature),
			Address:   ethAccRes.Address,
		})
		assert.NoError(s.T(), err)
	})
}
