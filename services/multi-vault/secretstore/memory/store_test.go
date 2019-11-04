package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/testutils"
)

type MemoryKeyStoreTestSuite struct {
	testutils.SecretStoreTestSuite
}

func (s *MemoryKeyStoreTestSuite) SetupTest() {
	s.Store = NewSecretStore()
}

func TestMemory(t *testing.T) {
	s := new(MemoryKeyStoreTestSuite)
	suite.Run(t, s)
}
