// +build integration

package integrationtests

import (
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	http2 "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/metrics"
	metrics1 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	tpcmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"gopkg.in/h2non/gock.v1"
)

type metricsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *metricsTestSuite) TestApplicationMetrics() {
	ctx := s.env.ctx
	defer gock.Off()

	chain := testutils.FakeChain()
	txRequest := testutils.FakeSendTransferTransactionRequest()
	gock.New(chainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
	gock.New(chainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
	gock.New(chainRegistryURL).
		URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", chainRegistryURL, chain.UUID, txRequest.Params.From)).
		Reply(404)

	s.T().Run("should increase created job metrics", func(t *testing.T) {
		mfsb, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)
		expectedV, err := testutils2.FamilyValue(mfsb, fmt.Sprintf("%s_%s", metrics1.Namespace, metrics.Subsystem), metrics.JobLatencySeconds, nil)
		if err != nil {
			expectedV = []uint64{0}
		}

		incrUintArr(expectedV.([]uint64))

		_, _ = s.client.SendTransferTransaction(ctx, txRequest)

		mfsa, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)
		testutils2.AssertFamilyValue(t, mfsa, fmt.Sprintf("%s_%s", metrics1.Namespace, metrics.Subsystem), metrics.JobLatencySeconds, expectedV, "", nil)
	})
}

func (s *metricsTestSuite) TestTCP() {
	ctx := s.env.ctx
	defer gock.Off()

	chain := testutils.FakeChain()
	txRequest := testutils.FakeSendTransferTransactionRequest()
	gock.New(chainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
	gock.New(chainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
	gock.New(chainRegistryURL).
		URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", chainRegistryURL, chain.UUID, txRequest.Params.From)).
		Reply(404)

	s.T().Run("should have one open entrypoint connection per type ('app', 'metrics')", func(t *testing.T) {
		mfs, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)

		testutils2.AssertFamilyValue(t, mfs, tpcmetrics.Namespace, tpcmetrics.OpenConns, []float64{1.0}, "open connections('app')", map[string]string{
			"entrypoint": "app",
		})
		testutils2.AssertFamilyValue(t, mfs, tpcmetrics.Namespace, tpcmetrics.OpenConns, []float64{1.0}, "open connections('metrics')", map[string]string{
			"entrypoint": "metrics",
		})
	})
}

func (s *metricsTestSuite) TestHTTP() {
	ctx := s.env.ctx
	defer gock.Off()

	chain := testutils.FakeChain()
	txRequest := testutils.FakeSendTransferTransactionRequest()
	gock.New(chainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
	gock.New(chainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
	gock.New(chainRegistryURL).
		URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", chainRegistryURL, chain.UUID, txRequest.Params.From)).
		Reply(404)

	s.T().Run("should successfully increase application total requests", func(t *testing.T) {
		mfsb, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)

		expectedV, err := testutils2.FamilyValue(mfsb, httpmetrics.Namespace, httpmetrics.RequestsTotal, map[string]string{
			"method": "POST",
		})
		if err != nil {
			expectedV = []float64{0.0}
		}

		incrFloatArr(expectedV.([]float64))

		_, err = s.client.SendTransferTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		time.Sleep(time.Second * 5)
		mfsa, err := s.client.Prometheus(ctx)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		testutils2.AssertFamilyValue(t, mfsa, httpmetrics.Namespace, httpmetrics.RequestsTotal, expectedV, "total POST requests", map[string]string{
			"method": "POST",
		})
	})
}

func (s *metricsTestSuite) TestZHealthCheck() {
	type healthRes struct {
		Database string `json:"database,omitempty"`
		Kafka    string `json:"kafka,omitempty"`
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	ctx := s.env.ctx
	s.T().Run("should retrieve positive health check over service dependencies", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

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
		assert.Equal(s.T(), "OK", status.Kafka)
	})

	s.T().Run("should retrieve a negative health check over kafka service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(chainRegistryMetricsURL).Get("/live").Reply(200)
		gock.New(contractRegistryMetricsURL).Get("/live").Reply(200)
		defer gock.Off()

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
	})

	s.T().Run("should retrieve a negative health check over postgres service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

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
	})
}

func incrFloatArr(arr []float64) {
	for i := range arr {
		arr[i]++
	}
}

func incrUintArr(arr []uint64) {
	for i := range arr {
		arr[i]++
	}
}
