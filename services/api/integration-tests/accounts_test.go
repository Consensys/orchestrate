// +build integration

package integrationtests

import (
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	http2 "net/http"
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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"gopkg.in/h2non/gock.v1"
)

type accountsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *accountsTestSuite) TestCreateAccounts() {
	ctx := s.env.ctx

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

		account := testutils.FakeETHAccountResponse()
		txRequest := testutils.FakeCreateAccountRequest()
		faucet := testutils.FakeFaucet()
		faucet.Creditor = ethcommon.HexToAddress("0x12278c8C089ef98b4045f0b649b61Ed4316B1a50")
		chain := testutils.FakeChain()
		txRequest.Chain = chain.Name

		// Create account and get faucet candidate for the newly created account
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(account)
		gock.New(chainRegistryURL).URL(fmt.Sprintf("%s/chains?name=%s", chainRegistryURL, chain.Name)).Times(2).Reply(200).JSON([]*models.Chain{chain})
		gock.New(chainRegistryURL).URL(fmt.Sprintf("%s/faucets/candidate?chain_uuid=%s&account=%s", chainRegistryURL, chain.UUID, account.Address)).Reply(200).JSON(faucet)

		// Send funding tx (this will check if the faucet account itself needs funding)
		gock.New(chainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chain)
		gock.New(chainRegistryURL).URL(fmt.Sprintf("%s/faucets/candidate?chain_uuid=%s&account=%s", chainRegistryURL, chain.UUID, faucet.Creditor.Hex())).Reply(404)

		resp, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Address, account.Address)
		assert.Equal(t, resp.PublicKey, account.PublicKey)
		assert.Equal(t, resp.Alias, txRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")
	})

	s.T().Run("should fail to create account if key-manager API fails", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(500).JSON(ethAccRes)

		_, err := s.client.CreateAccount(ctx, txRequest)
		assert.Error(s.T(), err)
	})

	s.T().Run("should fail to create account if postgres is down", func(t *testing.T) {
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		_, err = s.client.CreateAccount(ctx, txRequest)
		assert.Error(s.T(), err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(s.T(), err)
	})
}

func (s *accountsTestSuite) TestImportAccounts() {
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

		txRequest.Alias = fmt.Sprintf("Alias_%s", utils.RandomString(5))
		_, err = s.client.ImportAccount(ctx, txRequest)
		assert.Error(t, err)
		log.WithoutContext().Errorf("%v", err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})
}

func (s *accountsTestSuite) TestSearchIdentities() {
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

func (s *accountsTestSuite) TestGetAccount() {
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

func (s *accountsTestSuite) TestUpdateAccount() {
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
		gock.New(keyManagerURL).Post("/ethereum/accounts/verify-signature").Reply(http2.StatusNoContent)

		verifyRequest := testutils.FakeVerifyPayloadRequest()
		err := s.client.VerifySignature(ctx, verifyRequest)
		assert.NoError(t, err)
	})
}

func (s *accountsTestSuite) TestVerifyTypedDataSignature() {
	ctx := s.env.ctx

	s.T().Run("should verify typed data signature successfully", func(t *testing.T) {
		defer gock.Off()
		gock.New(keyManagerURL).Post("/ethereum/accounts/verify-typed-data-signature").Reply(http2.StatusNoContent)

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
