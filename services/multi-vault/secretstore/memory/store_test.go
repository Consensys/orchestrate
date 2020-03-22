// +build unit

package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/testutils"
)

// TODO: add new test with multi-tenancy context value

type MemoryKeyStoreTestSuite struct {
	testutils.SecretStoreTestSuite
}

func (s *MemoryKeyStoreTestSuite) SetupTest() {
	s.Store = NewSecretStore(multitenancy.New(false))
}

func TestMemory(t *testing.T) {
	s := new(MemoryKeyStoreTestSuite)
	suite.Run(t, s)
}
