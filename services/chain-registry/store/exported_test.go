package store

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mocks"
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

func (s *TestSuite) TestImportChains() {
	testSuite := []struct {
		name   string
		chains []string
	}{
		{
			"import chains",
			[]string{`{"name":"noError"}`},
		},
		{
			"import chains with error",
			[]string{`{"name":"error"}`},
		},
		{
			"import chains with unknown field",
			[]string{`{"unknown":"error"}`},
		},
	}

	mockCtrl := gomock.NewController(s.T())
	defer mockCtrl.Finish()
	mockStore := mocks.NewMockChainRegistryStore(mockCtrl)
	mockStore.EXPECT().RegisterChain(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, chain *types.Chain) error {
			if chain.Name == "error" {
				return fmt.Errorf("error")
			}
			return nil
		}).AnyTimes()
	mockStore.EXPECT().UpdateChainByName(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	for _, test := range testSuite {
		test := test
		s.T().Run(test.name, func(t *testing.T) {
			t.Parallel()

			importChains(context.Background(), test.chains, mockStore)
		})
	}

}
