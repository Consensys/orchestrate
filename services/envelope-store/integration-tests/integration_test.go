// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type envelopeStoreEnvTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func (s *envelopeStoreEnvTestSuite) SetupSuite() {
	s.env.Start()
}

func (s *envelopeStoreEnvTestSuite) TearDownSuite() {
	s.env.Teardown()
}

func TestEnvelopeStoreEnv_Init(t *testing.T) {
	s := new(envelopeStoreEnvTestSuite)
	s.env = NewIntegrationEnvironment(context.Background())
	suite.Run(t, s)
}

func (s *envelopeStoreEnvTestSuite) TestEnvelopeStore_GRPC() {
	grpcSuite := new(EnvelopeStoreTestSuite)
	grpcSuite.env = s.env
	grpcSuite.baseURL = "localhost:8080"
	suite.Run(s.T(), grpcSuite)
}

