package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/testutils"
)

type MockTestSuite struct {
	testutils.ContractRegistryTestSuite
	redisMock *miniredis.Miniredis
}

func (s *MockTestSuite) SetupTest() {
	redisMock, err := miniredis.Run()
	if err != nil {
		log.Fatalf("Could not start miniredis: %v", err.Error())
	}

	config := Config()
	config.URI = redisMock.Addr()

	s.R = NewRegistry(NewPool(config, Dial))
	s.redisMock = redisMock
}

func TestMock(t *testing.T) {
	s := new(MockTestSuite)
	suite.Run(t, s)
}
