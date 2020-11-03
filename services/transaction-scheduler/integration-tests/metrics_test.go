// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/metrics"
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

	s.T().Run("should capture application metrics", func(t *testing.T) {
		mfsb, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)
		expectedV, err := testutils2.FamilyValue(mfsb, metrics.MetricsNamespace, metrics.MetricCreatedJobName)
		if err != nil {
			expectedV = []float64{0.0}
		}

		for i := range expectedV.([]float64) {
			expectedV.([]float64)[i]++
		}

		_, err = s.client.SendTransferTransaction(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		mfsa, err := s.client.Prometheus(ctx)
		assert.NoError(t, err)
		testutils2.AssertFamilyValue(t, mfsa, metrics.MetricsNamespace, metrics.MetricCreatedJobName, expectedV, "")
	})
}
