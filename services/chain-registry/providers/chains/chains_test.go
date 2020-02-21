package chains

import (
	"context"
	"testing"
	"time"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/safe"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

var mockChainRegistryStore *mocks.MockChainRegistryStore

type ProviderTestSuite struct {
	suite.Suite
	provider *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	mockChainRegistryStore = mocks.NewMockChainRegistryStore(ctrl)

	s.provider = NewProvider(viper.GetString(store.TypeViperKey), mockChainRegistryStore, 200*time.Millisecond)
}

func (s *ProviderTestSuite) TestProvide() {
	mockChains := []*types.Chain{
		{
			Name:     "chain1",
			TenantID: "tenant1",
			URLs:     []string{"testUrl11", "testUrl12"},
		},
		{
			Name:     "chain2",
			TenantID: "tenant2",
			URLs:     []string{"testUrl21", "testUrl22"},
		},
	}

	mockChainRegistryStore.EXPECT().GetChains(gomock.Any(), gomock.Any()).Return(mockChains, nil).AnyTimes()

	assert.NoError(s.T(), s.provider.Init(), "Should initialize without error")

	ctx := context.Background()
	providerConfigUpdateCh := make(chan dynamic.Message)
	pool := safe.NewPool(ctx)

	go func() {
		err := s.provider.Provide(providerConfigUpdateCh, pool)
		assert.NoError(s.T(), err, "Should Provide without error")
	}()

	config := <-providerConfigUpdateCh
	assert.Equal(s.T(), viper.GetString(store.TypeViperKey), config.ProviderName, "Should get the correct providerName")
	close(providerConfigUpdateCh)
	pool.Stop()
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
