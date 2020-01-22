package store

import (
	"context"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg"
)

type TestSuite struct {
	suite.Suite
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	store = nil
	initOnce = &sync.Once{}
}

func (s *TestSuite) TestInit() {
	Init(context.Background())
	assert.NotNil(s.T(), GlobalStoreRegistry(), "Global should have been set")

	var chainRegistry *pg.ChainRegistry
	SetGlobalStoreRegistry(chainRegistry)
	assert.Nil(s.T(), GlobalStoreRegistry(), "Global should be reset to nil")
}

func (s *TestSuite) TestInitPostgres() {
	viper.Set(TypeViperKey, postgresOpt)
	Init(context.Background())
	assert.NotNil(s.T(), GlobalStoreRegistry(), "Global should have been set")

	var chainRegistry *pg.ChainRegistry
	SetGlobalStoreRegistry(chainRegistry)
	assert.Nil(s.T(), GlobalStoreRegistry(), "Global should be reset to nil")
}

func (s *TestSuite) TestInitInMemory() {
	viper.Set(TypeViperKey, memoryOpt)
	Init(context.Background())
	assert.NotNil(s.T(), GlobalStoreRegistry(), "Global should have been set")

	var chainRegistry *pg.ChainRegistry
	SetGlobalStoreRegistry(chainRegistry)
	assert.Nil(s.T(), GlobalStoreRegistry(), "Global should be reset to nil")
}

func (s *TestSuite) TestInitDefault() {
	defer func() { log.StandardLogger().ExitFunc = nil }()
	var fatal bool
	log.StandardLogger().ExitFunc = func(int) { fatal = true }

	viper.Set(TypeViperKey, "unknown")
	Init(context.Background())

	assert.True(s.T(), fatal, "should get fatal")

	var chainRegistry *pg.ChainRegistry
	SetGlobalStoreRegistry(chainRegistry)
	assert.Nil(s.T(), GlobalStoreRegistry(), "Global should be reset to nil")
}
