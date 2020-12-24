// +build integration

package integrationtests

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/integration-test"
)

type chainRegistryTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *chainRegistryTestSuite) SetupSuite() {
	err := integrationtest.StartEnvironment(s.env.ctx, s.env)
	if err != nil {
		s.env.logger.WithError(err).Error()
		if s.err == nil {
			s.err = err
		}
		return
	}

	s.env.logger.Infof("setup test suite has completed")
}

func (s *chainRegistryTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestChainRegistry(t *testing.T) {
	s := new(chainRegistryTestSuite)
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

func (s *chainRegistryTestSuite) TestChainRegistry_Chains() {
	httpSuite := new(chainsTestSuite)
	httpSuite.env = s.env
	httpSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), httpSuite)
}
