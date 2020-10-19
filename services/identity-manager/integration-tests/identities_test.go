// +build integration

package integrationtests

import (
	"fmt"
	http2 "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/client"
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

func (s *identityManagerTransactionTestSuite) TestTransactionScheduler_Transactions() {
	ctx := s.env.ctx
	ethAccRes := testutils.FakeETHAccountResponse()

	s.T().Run("should create identity successfully by querying key-manager API", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeCreateIdentityRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		resp, err := s.client.CreateIdentity(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Address, ethAccRes.Address)
		assert.Equal(t, resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(t, resp.Alias, txRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")
	})

	s.T().Run("should import identity successfully by querying key-manager API", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeImportIdentityRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts/import").Reply(200).JSON(ethAccRes)

		resp, err := s.client.ImportIdentity(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Address, ethAccRes.Address)
		assert.Equal(t, resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(t, resp.Alias, txRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")
	})

	s.T().Run("should create identity successfully and trigger funding transaction", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeCreateIdentityRequest()
		txRequest.Chain = "besu"
		faucet := testutils.FakeFaucet()
		chain := testutils.FakeChain()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)
		gock.New(chainRegistryURL).Post(fmt.Sprintf("/chains?name=%s", txRequest.Chain)).Reply(200).JSON(chain)
		gock.New(chainRegistryURL).Post(fmt.Sprintf("/faucets/candidate?chain_uuid=%s&account=%s", chain.UUID, ethAccRes.Address)).Reply(200).JSON(faucet)
		gock.New(chainRegistryURL).Post("/transactions/transfer").Reply(200)

		resp, err := s.client.CreateIdentity(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, resp.Address, ethAccRes.Address)
		assert.Equal(t, resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(t, resp.Alias, txRequest.Alias)
		assert.Equal(t, resp.TenantID, "_")
	})

	s.T().Run("should fail to create identity if key-manager API fails", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeCreateIdentityRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(500).JSON(ethAccRes)

		_, err := s.client.CreateIdentity(ctx, txRequest)
		assert.Error(s.T(), err)
	})

	s.T().Run("should fail to create identity if postgres is down", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeCreateIdentityRequest()
		gock.New(keyManagerURL).Post("/ethereum/accounts").Reply(200).JSON(ethAccRes)

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		_, err = s.client.CreateIdentity(ctx, txRequest)
		assert.Error(s.T(), err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(s.T(), err)
	})
}

func (s *identityManagerTransactionTestSuite) TestTransactionScheduler_ZHealthCheck() {
	type healthRes struct {
		KeyManager    string `json:"key-manager,omitempty"`
		ChainRegistry string `json:"chain-registry,omitempty"`
		TxScheduler   string `json:"transaction-scheduler,omitempty"`
		Database      string `json:"database,omitempty"`
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	ctx := s.env.ctx

	s.T().Run("should retrieve positive health check over service dependencies", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(200)
		gock.New(txSchedulerMetricsURL).Get("/live").Reply(200)
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
		assert.Equal(s.T(), "OK", status.TxScheduler)
	})

	s.T().Run("should retrieve a negative health check over key-manager API service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(500)
		gock.New(txSchedulerMetricsURL).Get("/live").Reply(200)
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
		assert.Equal(s.T(), "OK", status.TxScheduler)
	})

	s.T().Run("should retrieve a negative health check over chain-registry API service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(200)
		gock.New(txSchedulerMetricsURL).Get("/live").Reply(200)
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
		assert.Equal(s.T(), "OK", status.TxScheduler)
	})

	s.T().Run("should retrieve a negative health check over tx-scheduler API service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(200)
		gock.New(txSchedulerMetricsURL).Get("/live").Reply(500)
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
		assert.Equal(s.T(), "OK", status.KeyManager)
		assert.Equal(s.T(), "OK", status.ChainRegistry)
		assert.NotEqual(s.T(), "OK", status.TxScheduler)
	})

	s.T().Run("should retrieve a negative health check over postgres service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(keyManagerMetricsURL).Get("/live").Reply(200)
		gock.New(txSchedulerMetricsURL).Get("/live").Reply(200)
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
		assert.Equal(s.T(), "OK", status.TxScheduler)
	})
}
