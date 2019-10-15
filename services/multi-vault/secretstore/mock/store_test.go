package mock

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/multi-vault/secretstore/testutils"
)

type MockKeyStoreTestSuite struct {
	testutils.SecretStoreTestSuite
}

func (s *MockKeyStoreTestSuite) SetupTest() {
	s.Store = NewSecretStore()
}

func TestMock(t *testing.T) {
	s := new(MockKeyStoreTestSuite)
	suite.Run(t, s)
}