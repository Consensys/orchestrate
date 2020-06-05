// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
)

type contractRegistryTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *contractRegistryTestSuite) SetupSuite() {
	err := integrationtest.StartEnvironment(s.env)
	if err != nil {
		s.env.logger.WithError(err).Error()
		if s.err == nil {
			s.err = err
		}
		return
	}

	s.env.logger.Infof("setup test suite has completed")
}

func (s *contractRegistryTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestContractRegistry(t *testing.T) {
	s := new(contractRegistryTestSuite)
	s.env, s.err = NewIntegrationEnvironment(context.Background())
	if s.err != nil {
		t.Fail()
		return
	}

	suite.Run(t, s)
}

func (s *contractRegistryTestSuite) TestContractRegistry_HTTP() {
	if s.err != nil {
		s.env.logger.Warn("skipping test TestContractRegistry_HTTP...")
		return
	}
	
	httpSuite := new(contractRegistryHTTPTestSuite)
	httpSuite.env = s.env
	httpSuite.baseURL = s.env.baseHTTP
	suite.Run(s.T(), httpSuite)
}

func (s *contractRegistryTestSuite) TestContractRegistry_GRPC() {
	if s.err != nil {
		s.env.logger.Warn("skipping test TestContractRegistry_GRPC...")
		return
	}

	grpcSuite := new(contractRegistryGRPCTestSuite)
	grpcSuite.env = s.env
	grpcSuite.baseURL = s.env.baseGRPC
	suite.Run(s.T(), grpcSuite)
}
