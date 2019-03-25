package mock

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/testutils"
)

type MockKeyStoreTestSuite struct {
	testutils.SecretStoreTestSuite
}

func (suite *MockKeyStoreTestSuite) SetupTest() {
	suite.Store = NewSecretStore()
}

func TestMock(t *testing.T) {
	s := new(MockKeyStoreTestSuite)
	suite.Run(t, s)
}
