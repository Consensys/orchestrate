// +build integration

package integrationtests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type contractRegistryTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func TestContractRegistry(t *testing.T) {
	s := new(contractRegistryTestSuite)
	s.env = NewIntegrationEnvironment()
	suite.Run(t, s)
}

func (s *contractRegistryTestSuite) SetupSuite() {
	s.env.Start()
}

func (s *contractRegistryTestSuite) TearDownSuite() {
	s.env.Teardown()
}

func (s *contractRegistryTestSuite) TestContractRegistry_HTTP() {
	httpSuite := new(contractRegistryHTTPTestSuite)
	httpSuite.env = s.env
	httpSuite.baseURL = "http://localhost:8081"
	suite.Run(s.T(), httpSuite)
}

func (s *contractRegistryTestSuite) TestContractRegistry_GRPC() {
	grpcSuite := new(contractRegistryGRPCTestSuite)
	grpcSuite.env = s.env
	grpcSuite.baseURL = "localhost:8080"
	suite.Run(s.T(), grpcSuite)
}
