// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type envelopeStoreTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func (s *envelopeStoreTestSuite) SetupSuite() {
	s.env.Start()
}

func (s *envelopeStoreTestSuite) TearDownSuite() {
	s.env.Teardown()
}

func TestEnvelopeStore(t *testing.T) {
	s := new(envelopeStoreTestSuite)
	s.env = NewIntegrationEnvironment(context.Background())
	suite.Run(t, s)
}

func (s *envelopeStoreTestSuite) TestEnvelopeStore_GRPC() {
	grpcSuite := new(EnvelopeStoreTestSuite)
	grpcSuite.env = s.env
	grpcSuite.baseURL = "localhost:8080"
	suite.Run(s.T(), grpcSuite)
}

