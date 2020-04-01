// +build unit

package chainregistry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

var mockChainRegistryClient *mocks.MockChainRegistryClient

type ProviderTestSuite struct {
	suite.Suite
	provider *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	mockChainRegistryClient = mocks.NewMockChainRegistryClient(ctrl)

	s.provider = &Provider{
		Client: mockChainRegistryClient,
		conf: &Config{
			RefreshInterval:  time.Millisecond,
			ChainRegistryURL: "http://test-proxy",
		},
	}
}

func (s *ProviderTestSuite) TestRun() {
	mockChains := []*types.Chain{
		{
			UUID:                      "0d60a85e-0b90-4482-a14c-108aea2557aa",
			Name:                      "chainName",
			TenantID:                  "0d60a85e-0b90-4482-a14c-108aea2557bb",
			URLs:                      []string{"https://estcequecestbientotlapero.fr/"},
			ListenerDepth:             &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:      &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:     &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration:   &(&struct{ x string }{"1s"}).x,
			ListenerExternalTxEnabled: &(&struct{ x bool }{true}).x,
		},
	}

	gomock.InOrder(
		mockChainRegistryClient.EXPECT().GetChains(gomock.Any()).Return([]*types.Chain{}, nil),
		mockChainRegistryClient.EXPECT().GetChains(gomock.Any()).Return(mockChains, nil).AnyTimes(),
		mockChainRegistryClient.EXPECT().GetChains(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes(),
	)

	cancelableCtx, cancel := context.WithCancel(context.Background())
	providerConfigUpdateCh := make(chan *dynamic.Message)
	go func() {
		runErr := s.provider.Run(cancelableCtx, providerConfigUpdateCh)
		assert.NoError(s.T(), runErr)
		close(providerConfigUpdateCh)
	}()
	config := <-providerConfigUpdateCh
	assert.Equal(s.T(), "chain-registry", config.Provider, "Should get the correct providerName")
	assert.Len(s.T(), config.Configuration.Chains, 0)

	config = <-providerConfigUpdateCh
	assert.Equal(s.T(), "chain-registry", config.Provider, "Should get the correct providerName")
	assert.Len(s.T(), config.Configuration.Chains, 1)
	assert.Equal(
		s.T(),
		"http://test-proxy/0d60a85e-0b90-4482-a14c-108aea2557aa",
		config.Configuration.Chains["0d60a85e-0b90-4482-a14c-108aea2557aa"].URL,
		"Chain URL should be correct",
	)

	cancel()
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
