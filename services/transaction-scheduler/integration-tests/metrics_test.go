// +build integration

package integrationtests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/metrics"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/testutils"
	tpcmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/identitymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/metrics"
	"gopkg.in/h2non/gock.v1"
)

type txSchedulerMetricsTestSuite struct {
	suite.Suite
	baseURL    string
	metricsURL string
	client     client.TransactionSchedulerClient
	env        *IntegrationEnvironment
}

func (s *txSchedulerMetricsTestSuite) SetupSuite() {
	conf := client.NewConfig(s.baseURL, nil)
	conf.MetricsURL = s.metricsURL
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)
}

func (s *txSchedulerMetricsTestSuite) TestTransactionMetrics_Application() {
	ctx := s.env.ctx
	defer gock.Off()

	chain := testutils.FakeChain()
	txRequest := testutils.FakeSendTransferTransactionRequest()
	gock.New(chainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
	gock.New(chainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
	gock.New(identityManagerURL).Get("/accounts/" + txRequest.Params.From).Reply(200).JSON(&identitymanager.AccountResponse{})
	gock.New(chainRegistryURL).
		URL(fmt.Sprintf("%s?chain_uuid=%s&account=%s", chainRegistryURL, chain.UUID, txRequest.Params.From)).
		Reply(404)

	s.T().Run("should increase created job metrics", func(t *testing.T) {
		mfsb, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)
		expectedV, err := testutils2.FamilyValue(mfsb, metrics.Namespace, metrics.CreatedJobName, nil)
		if err != nil {
			expectedV = []float64{0.0}
		}

		incrFloatArr(expectedV.([]float64))

		_, _ = s.client.SendTransferTransaction(ctx, txRequest)

		mfsa, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)
		testutils2.AssertFamilyValue(t, mfsa, metrics.Namespace, metrics.CreatedJobName, expectedV, "", nil)
	})
}

func (s *txSchedulerMetricsTestSuite) TestTransactionMetrics_TCP() {
	ctx := s.env.ctx
	defer gock.Off()

	chain := testutils.FakeChain()
	txRequest := testutils.FakeSendTransferTransactionRequest()
	gock.New(chainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
	gock.New(chainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
	gock.New(identityManagerURL).Get("/accounts/" + txRequest.Params.From).Reply(200).JSON(&identitymanager.AccountResponse{})
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

func (s *txSchedulerMetricsTestSuite) TestTransactionMetrics_HTTP() {
	ctx := s.env.ctx
	defer gock.Off()

	chain := testutils.FakeChain()
	txRequest := testutils.FakeSendTransferTransactionRequest()
	gock.New(chainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chain})
	gock.New(chainRegistryURL).Get("/chains/" + chain.UUID).Times(2).Reply(200).JSON(chain)
	gock.New(identityManagerURL).Get("/accounts/" + txRequest.Params.From).Reply(200).JSON(&identitymanager.AccountResponse{})
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

		_, _ = s.client.SendTransferTransaction(ctx, txRequest)
		time.Sleep(time.Second * 5)
		mfsa, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)

		testutils2.AssertFamilyValue(t, mfsa, httpmetrics.Namespace, httpmetrics.RequestsTotal, expectedV, "total POST requests", map[string]string{
			"method": "POST",
		})
	})
}

func incrFloatArr(arr []float64) {
	for i := range arr {
		arr[i]++
	}
}
