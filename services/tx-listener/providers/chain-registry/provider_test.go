package chainregistry

import (
	"context"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type MockClient struct {
	i int
}

func (m *MockClient) GetNodeByID(_ context.Context, _ string) (*types.Node, error) {
	return &types.Node{}, nil
}

func (m *MockClient) GetNodeByTenantAndNodeName(_ context.Context, _, _ string) (*types.Node, error) {
	return nil, nil
}

func (m *MockClient) GetNodeByTenantAndNodeID(_ context.Context, _, _ string) (*types.Node, error) {
	return nil, nil
}

func (m *MockClient) GetNodes(_ context.Context) ([]*types.Node, error) {
	switch m.i % 2 {
	case 0:
		m.i++
		return []*types.Node{}, nil
	case 1:
		m.i++
		return []*types.Node{
			{
				ID:       "0d60a85e-0b90-4482-a14c-108aea2557aa",
				Name:     "42",
				TenantID: "0d60a85e-0b90-4482-a14c-108aea2557bb",
				URLs:     []string{"https://estcequecestbientotlapero.fr/"},
			}}, nil
	default:
		return []*types.Node{}, nil
	}
}

func (m *MockClient) UpdateBlockPosition(_ context.Context, _ string, _ int64) error {
	return nil
}

type ProviderTestSuite struct {
	suite.Suite
	provider *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	s.provider = &Provider{
		Client: &MockClient{i: 0},
		conf: &Config{
			RefreshInterval:  time.Millisecond,
			ChainRegistryURL: "http://test-proxy",
		},
	}
}

func (s *ProviderTestSuite) TestRun() {
	cancelableCtx, cancel := context.WithCancel(context.Background())
	providerConfigUpdateCh := make(chan *dynamic.Message)
	go func() {
		runErr := s.provider.Run(cancelableCtx, providerConfigUpdateCh)
		assert.NoError(s.T(), runErr)
		close(providerConfigUpdateCh)
	}()
	config := <-providerConfigUpdateCh
	assert.Equal(s.T(), "chain-registry", config.Provider, "Should get the correct providerName")
	assert.Len(s.T(), config.Configuration.Nodes, 0)

	config = <-providerConfigUpdateCh
	assert.Equal(s.T(), "chain-registry", config.Provider, "Should get the correct providerName")
	assert.Len(s.T(), config.Configuration.Nodes, 1)
	assert.Equal(
		s.T(),
		"http://test-proxy/0d60a85e-0b90-4482-a14c-108aea2557aa",
		config.Configuration.Nodes["0d60a85e-0b90-4482-a14c-108aea2557aa"].URL,
		"Node URL should be correct",
	)

	cancel()
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
