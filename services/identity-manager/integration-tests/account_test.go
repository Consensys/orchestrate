// +build integration

package integrationtests

import (
	"fmt"
	http2 "net/http"
	"testing"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/identitymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/client"
	"gopkg.in/h2non/gock.v1"
)

type identityManagerTransactionTestSuite struct {
	suite.Suite
	baseURL string
	client  client.IdentityManagerClient
	env     *IntegrationEnvironment
}

func (s *identityManagerTransactionTestSuite) SetupSuite() {
	conf := client.NewConfig(s.baseURL, nil)
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)
}

func (s *identityManagerTransactionTestSuite) TestIdentityManager_CreateAccounts() {
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
		ethAccRes := testutils.FakeETHAccountResponse()
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()
		faucet := testutils.FakeFaucet()
		chain := testutils.FakeChain()
		txRequest.Chain = chain.Name
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)
		gock.New(chainRegistryURL).URL(fmt.Sprintf("%s/chains?name=%s", chainRegistryURL, chain.Name)).Reply(200).JSON([]*models.Chain{chain})
		gock.New(chainRegistryURL).URL(fmt.Sprintf("%s/faucets/candidate?chain_uuid=%s&account=%s", chainRegistryURL, chain.UUID, ethAccRes.Address)).Reply(200).JSON(faucet)
		gock.New(txSchedulerURL).Post("/transactions/transfer").Reply(200).JSON(txscheduler.TransactionResponse{})

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

func (s *identityManagerTransactionTestSuite) TestIdentityManager_ImportAccounts() {
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

func (s *identityManagerTransactionTestSuite) TestIdentityManager_SearchIdentities() {
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

func (s *identityManagerTransactionTestSuite) TestIdentityManager_GetAccount() {
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

func (s *identityManagerTransactionTestSuite) TestIdentityManager_UpdateAccount() {
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

func (s *identityManagerTransactionTestSuite) TestIdentityManager_SignPayload() {
	ctx := s.env.ctx

	s.T().Run("should sign payload successfully", func(t *testing.T) {
		defer gock.Off()
		address := ethcommon.HexToAddress("0x123").String()
		payload := "messageToSign"
		signedPayload := ethcommon.HexToHash("0xABCDEF01234").String()
		gock.New(keyManagerURL).Post(fmt.Sprintf("/ethereum/accounts/%s/sign", address)).
			Reply(200).BodyString(signedPayload)

		response, err := s.client.SignPayload(ctx, address, &types.SignPayloadRequest{
			Data: payload,
		})
		assert.NoError(t, err)
		assert.Equal(t, signedPayload, response)
	})
}

func (s *identityManagerTransactionTestSuite) TestIdentityManager_VerifySignature() {
	ctx := s.env.ctx

	s.T().Run("should verify signature successfully", func(t *testing.T) {
		defer gock.Off()
		gock.New(keyManagerURL).Post("/ethereum/accounts/verify-signature").Reply(http2.StatusNoContent)

		verifyRequest := testutils.FakeVerifyPayloadRequest()
		err := s.client.VerifySignature(ctx, verifyRequest)
		assert.NoError(t, err)
	})
}

func (s *identityManagerTransactionTestSuite) TestIdentityManager_VerifyTypedDataSignature() {
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

func (s *identityManagerTransactionTestSuite) TestIdentityManager_SignTypedData() {
	ctx := s.env.ctx

	s.T().Run("should sign typed data successfully", func(t *testing.T) {
		defer gock.Off()
		address := ethcommon.HexToAddress("0x123").String()
		signature := "0xsignature"
		gock.New(keyManagerURL).Post(fmt.Sprintf("/ethereum/accounts/%s/sign-typed-data", address)).
			Reply(200).BodyString(signature)

		typedDataRequest := testutils.FakeSignTypedDataRequest()
		response, err := s.client.SignTypedData(ctx, address, &types.SignTypedDataRequest{
			DomainSeparator: typedDataRequest.DomainSeparator,
			Types:           typedDataRequest.Types,
			Message:         typedDataRequest.Message,
			MessageType:     typedDataRequest.MessageType,
		})

		assert.NoError(t, err)
		assert.Equal(t, signature, response)
	})
}

func (s *identityManagerTransactionTestSuite) TestIdentityManager_ZHealthCheck() {
	type healthRes struct {
		KeyManager    string `json:"key-manager,omitempty"`
		ChainRegistry string `json:"chain-registry,omitempty"`
		Database      string `json:"database,omitempty"`
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	ctx := s.env.ctx

	s.T().Run("should retrieve positive health check over service dependencies", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(200)
		gock.New(chainRegistryMetricsURL).Get("/live").Reply(200)
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
		assert.Equal(s.T(), "OK", status.KeyManager)
		assert.Equal(s.T(), "OK", status.ChainRegistry)
	})

	s.T().Run("should retrieve a negative health check over key-manager API service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(500)
		gock.New(chainRegistryMetricsURL).Get("/live").Reply(200)
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
		assert.NotEqual(s.T(), "OK", status.KeyManager)
		assert.Equal(s.T(), "OK", status.ChainRegistry)
	})

	s.T().Run("should retrieve a negative health check over chain-registry API service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(200)
		gock.New(chainRegistryMetricsURL).Get("/live").Reply(500)
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
		assert.Equal(s.T(), "OK", status.KeyManager)
		assert.NotEqual(s.T(), "OK", status.ChainRegistry)
	})

	s.T().Run("should retrieve a negative health check over postgres service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(200)
		gock.New(chainRegistryMetricsURL).Get("/live").Reply(200)
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
		assert.Equal(s.T(), "OK", status.KeyManager)
		assert.Equal(s.T(), "OK", status.ChainRegistry)
	})
}
