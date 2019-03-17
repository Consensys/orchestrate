package mock

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/testutils"
)

type MockTraceStoreTestSuite struct {
	testutils.TraceStoreTestSuite
}

func (suite *MockTraceStoreTestSuite) SetupTest() {
	suite.Store = NewTraceStore()
}

func TestMock(t *testing.T) {
	s := new(MockTraceStoreTestSuite)
	suite.Run(t, s)
}
