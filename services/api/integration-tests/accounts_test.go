// +build integration

package integrationtests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"net/http"
	"testing"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gopkg.in/h2non/gock.v1"
)

type accountsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *accountsTestSuite) TestCreateAccounts() {
	ctx := s.env.ctx
	chain := testutils.FakeChain()
	chain.URLs = []string{s.env.blockchainNodeURL}

	s.T().Run("should create account successfully by querying key-manager API", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		resp, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Address, ethAccRes.Address)
		assert.Equal(t, resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(t, resp.Alias, txRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")
	})

	s.T().Run("should fail to create account with same alias", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Times(2).Reply(200).JSON(ethAccRes)

		_, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		_, err = s.client.CreateAccount(ctx, txRequest)
		assert.Error(t, err)
		log.WithoutContext().Errorf("%v", err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	s.T().Run("should create account successfully and trigger funding transaction", func(t *testing.T) {
		defer gock.Off()

		chainWithFaucet, err := s.client.RegisterChain(s.env.ctx, &api.RegisterChainRequest{
			Name: "ganache-with-faucet-accounts",
			URLs: []string{s.env.blockchainNodeURL},
			Listener: api.RegisterListenerRequest{
				FromBlock:         "latest",
				ExternalTxEnabled: false,
			},
		})
		require.NoError(t, err)

		accountFaucet := testutils.FakeAccount()
		accountFaucet.Alias = "MyFaucetCreditor-accounts"
		accountFaucet.Address = "0xcE187877Afa6C3830342958E1D9ab6E707e8f863"
		gock.New(keyManagerURL).Post("/ethereum/accounts/import").Reply(200).JSON(accountFaucet)
		_, err = s.client.ImportAccount(s.env.ctx, testutils.FakeImportAccountRequest())
		require.NoError(t, err)

		faucetRequest := testutils.FakeRegisterFaucetRequest()
		faucetRequest.Name = "faucet-integration-tests"
		faucetRequest.ChainRule = chainWithFaucet.UUID
		faucetRequest.CreditorAccount = accountFaucet.Address
		faucet, err := s.client.RegisterFaucet(s.env.ctx, faucetRequest)
		require.NoError(s.T(), err)

		account := testutils.FakeETHAccountResponse()
		accountRequest := testutils.FakeCreateAccountRequest()
		accountRequest.Chain = chainWithFaucet.Name

		// Create account and get faucet candidate for the newly created account
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(account)

		resp, err := s.client.CreateAccount(ctx, accountRequest)
		require.NoError(t, err)

		assert.Equal(t, resp.Address, account.Address)
		assert.Equal(t, resp.PublicKey, account.PublicKey)
		assert.Equal(t, resp.Alias, accountRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")

		err = s.client.DeleteChain(ctx, chainWithFaucet.UUID)
		assert.NoError(t, err)
		err = s.client.DeleteFaucet(ctx, faucet.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should fail to create account if key-manager API fails", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(500).JSON(ethAccRes)

		_, err := s.client.CreateAccount(ctx, txRequest)
		assert.Error(t, err)
	})

	s.T().Run("should fail to create account if postgres is down", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		_, err = s.client.CreateAccount(ctx, txRequest)
		assert.Error(t, err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(t, err)
	})
}

func (s *accountsTestSuite) TestImport() {
	ctx := s.env.ctx

	s.T().Run("should import account successfully by querying key-manager API", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeImportAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts/import").Reply(200).JSON(ethAccRes)

		resp, err := s.client.ImportAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Address, ethAccRes.Address)
		assert.Equal(t, resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(t, resp.Alias, txRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")
	})

	s.T().Run("should fail to import same account twice", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeImportAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts/import").Times(2).Reply(200).JSON(ethAccRes)

		_, err := s.client.ImportAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		txRequest.Alias = fmt.Sprintf("Alias_%s", utils.RandString(5))
		_, err = s.client.ImportAccount(ctx, txRequest)
		assert.Error(t, err)
		log.WithoutContext().Errorf("%v", err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})
}

func (s *accountsTestSuite) TestSearch() {
	ctx := s.env.ctx

	s.T().Run("should create account and search for it by alias successfully", func(t *testing.T) {
		defer gock.Off()
		ethAccRes := testutils.FakeETHAccountResponse()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		_, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		resp, err := s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: []string{txRequest.Alias},
		})
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Len(t, resp, 1)
		assert.Equal(t, resp[0].Address, ethAccRes.Address)
		assert.Equal(t, resp[0].PublicKey, ethAccRes.PublicKey)
		assert.Equal(t, resp[0].Alias, txRequest.Alias)
		assert.Equal(t, resp[0].TenantID, "_")
	})
}

func (s *accountsTestSuite) TestGetOne() {
	ctx := s.env.ctx

	s.T().Run("should create account and get it by address successfully", func(t *testing.T) {
		defer gock.Off()
		ethAccRes := testutils.FakeETHAccountResponse()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		_, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		resp, err := s.client.GetAccount(ctx, ethAccRes.Address)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Address, ethAccRes.Address)
		assert.Equal(t, resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(t, resp.Alias, txRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")
	})
}

func (s *accountsTestSuite) TestUpdate() {
	ctx := s.env.ctx

	s.T().Run("should create account and update it successfully", func(t *testing.T) {
		defer gock.Off()
		ethAccRes := testutils.FakeETHAccountResponse()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		_, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		txRequest2 := testutils.FakeUpdateAccountRequest()
		resp, err := s.client.UpdateAccount(ctx, ethAccRes.Address, txRequest2)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Alias, txRequest2.Alias)
		assert.Equal(t, resp.Attributes, txRequest2.Attributes)
		assert.Equal(t, resp.TenantID, "_")
	})
}

func (s *accountsTestSuite) TestSignPayload() {
	ctx := s.env.ctx

	s.T().Run("should sign payload successfully", func(t *testing.T) {
		defer gock.Off()
		address := ethcommon.HexToAddress("0x123").String()
		payload := "messageToSign"
		signedPayload := ethcommon.HexToHash("0xABCDEF01234").String()
		gock.New(keyManagerURL).Post(fmt.Sprintf("/ethereum/accounts/%s/sign", address)).
			Reply(200).BodyString(signedPayload)

		response, err := s.client.SignPayload(ctx, address, &api.SignPayloadRequest{
			Data: payload,
		})
		assert.NoError(t, err)
		assert.Equal(t, signedPayload, response)
	})
}

func (s *accountsTestSuite) TestVerifySignature() {
	ctx := s.env.ctx

	s.T().Run("should verify signature successfully", func(t *testing.T) {
		defer gock.Off()
		gock.New(keyManagerURL).Post("/ethereum/accounts/verify-signature").Reply(http.StatusNoContent)

		verifyRequest := testutils.FakeVerifyPayloadRequest()
		err := s.client.VerifySignature(ctx, verifyRequest)
		assert.NoError(t, err)
	})
}

func (s *accountsTestSuite) TestVerifyTypedDataSignature() {
	ctx := s.env.ctx

	s.T().Run("should verify typed data signature successfully", func(t *testing.T) {
		defer gock.Off()
		gock.New(keyManagerURL).Post("/ethereum/accounts/verify-typed-data-signature").Reply(http.StatusNoContent)

		verifyRequest := testutils.FakeVerifyTypedDataPayloadRequest()
		err := s.client.VerifyTypedDataSignature(ctx, verifyRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	})
}

func (s *accountsTestSuite) TestSignTypedData() {
	ctx := s.env.ctx

	s.T().Run("should sign typed data successfully", func(t *testing.T) {
		defer gock.Off()
		address := ethcommon.HexToAddress("0x123").String()
		signature := "0xsignature"
		gock.New(keyManagerURL).Post(fmt.Sprintf("/ethereum/accounts/%s/sign-typed-data", address)).
			Reply(200).BodyString(signature)

		typedDataRequest := testutils.FakeSignTypedDataRequest()
		response, err := s.client.SignTypedData(ctx, address, &api.SignTypedDataRequest{
			DomainSeparator: typedDataRequest.DomainSeparator,
			Types:           typedDataRequest.Types,
			Message:         typedDataRequest.Message,
			MessageType:     typedDataRequest.MessageType,
		})

		assert.NoError(t, err)
		assert.Equal(t, signature, response)
	})
}
