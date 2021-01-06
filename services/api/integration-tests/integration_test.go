// +build integration

package integrationtests

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gopkg.in/h2non/gock.v1"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/integration-test"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type apiTestSuite struct {
	suite.Suite
	env    *IntegrationEnvironment
	client client.OrchestrateClient
	faucet *api.FaucetResponse
	err    error
}

func (s *apiTestSuite) SetupSuite() {
	defer gock.Off()

	err := integrationtest.StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.env.logger.WithError(err).Error()
		if s.err == nil {
			s.err = err
		}
		return
	}

	s.env.logger.Debug("setting up test accounts and chains")

	conf := client.NewConfig(s.env.baseURL, nil)
	conf.MetricsURL = s.env.metricsURL
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)

	// We use this faucet in the tests
	faucetRequest := testutils.FakeRegisterFaucetRequest()
	faucetRequest.Name = "faucet-integration-tests"
	accountFaucet := testutils.FakeAccount()
	accountFaucet.Alias = "MyFaucetCreditor"
	accountFaucet.Address = faucetRequest.CreditorAccount
	gock.New(keyManagerURL).Post("/ethereum/accounts/import").Reply(200).JSON(accountFaucet)
	_, s.err = s.client.ImportAccount(s.env.ctx, testutils.FakeImportAccountRequest())
	if s.err != nil {
		s.T().Errorf(s.err.Error())
		return
	}
	s.faucet, s.err = s.client.RegisterFaucet(s.env.ctx, faucetRequest)
	if s.err != nil {
		s.T().Errorf(s.err.Error())
		return
	}

	// We use this account in the tests
	account := testutils.FakeAccount()
	account.Address = "0x5Cc634233E4a454d47aACd9fC68801482Fb02610"
	gock.New(keyManagerURL).Post("/ethereum/accounts/import").Reply(200).JSON(account)
	_, s.err = s.client.ImportAccount(s.env.ctx, testutils.FakeImportAccountRequest())
	if s.err != nil {
		s.T().Errorf(s.err.Error())
		return
	}

	s.env.logger.Infof("setup test suite has completed")
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

	time.Sleep(2 * time.Second)
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
	testSuite.faucet = s.faucet
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

/*
func (s *apiTestSuite) TestAPI_Metrics() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(metricsTestSuite)
	testSuite.env = s.env
	testSuite.client = s.client
	suite.Run(s.T(), testSuite)
}
*/
