// +build integration

package integrationtests

import (
	"context"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	integrationtest "github.com/consensys/orchestrate/pkg/toolkit/integration-test"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"os"
	"testing"
	"time"
)

type apiTestSuite struct {
	suite.Suite
	env       *IntegrationEnvironment
	client    client.OrchestrateClient
	chainUUID string
	err       error
}

func (s *apiTestSuite) SetupSuite() {
	defer gock.Off()

	s.err = integrationtest.StartEnvironment(s.env.ctx, s.env)
	if s.err != nil {
		s.Fail(s.err.Error())
	}
	time.Sleep(2 * time.Second)

	s.env.logger.Debug("setting up test accounts and chains")

	conf := client.NewConfig(s.env.baseURL, nil)
	conf.MetricsURL = s.env.metricsURL
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)

	// We use this chain in the tests
	chain, err := s.client.RegisterChain(s.env.ctx, &api.RegisterChainRequest{
		Name: "ganache",
		URLs: []string{s.env.blockchainNodeURL},
		Listener: api.RegisterListenerRequest{
			FromBlock:         "latest",
		},
	})
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.err = err
	s.chainUUID = chain.UUID

	// We use this account in the tests
	account := testutils.FakeAccount()
	account.Address = testutils.FromAddress
	_, s.err = s.client.ImportAccount(s.env.ctx, testutils.FakeImportAccountRequest())
	if s.err != nil {
		s.Fail(s.err.Error())
	}

	s.env.logger.Info("setup test suite has completed")
}

func (s *apiTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestAPI(t *testing.T) {
	s := new(apiTestSuite)
	ctx, cancel := context.WithCancel(context.Background())

	s.env, s.err = NewIntegrationEnvironment(ctx)
	if s.err != nil {
		t.Errorf(s.err.Error())
		return
	}

	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	suite.Run(t, s)
}

func (s *apiTestSuite) TestAPI_Transactions() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(transactionsTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	suite.Run(s.T(), testSuite)
}

func (s *apiTestSuite) TestAPI_Accounts() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(accountsTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	testSuite.defaultQKMStoreID = qkmDefaultStoreID
	suite.Run(s.T(), testSuite)
}

func (s *apiTestSuite) TestAPI_Jobs() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(jobsTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	testSuite.chainUUID = s.chainUUID
	suite.Run(s.T(), testSuite)
}

func (s *apiTestSuite) TestAPI_Schedules() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(schedulesTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	suite.Run(s.T(), testSuite)
}

func (s *apiTestSuite) TestAPI_Contracts() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(contractsTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	suite.Run(s.T(), testSuite)
}

func (s *apiTestSuite) TestAPI_Chains() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(chainsTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	suite.Run(s.T(), testSuite)
}

func (s *apiTestSuite) TestAPI_Faucets() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(faucetsTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	suite.Run(s.T(), testSuite)
}

func (s *apiTestSuite) TestAPI_Proxy() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(proxyTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	suite.Run(s.T(), testSuite)
}

// func (s *apiTestSuite) TestAPI_Metrics() {
// 	if s.err != nil {
// 		s.env.logger.Warn("skipping test...")
// 		return
// 	}
//
// 	testSuite := new(metricsTestSuite)
// 	testSuite.env = s.env
// 	testSuite.client = s.client
// 	suite.Run(s.T(), testSuite)
// }
