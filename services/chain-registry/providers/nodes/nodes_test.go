package nodes

import (
	"context"
	"testing"
	"time"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/safe"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/memory"
)

type ProviderTestSuite struct {
	suite.Suite
	provider *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	// TODO: create specific ChainRegistry for test instead of using memory chain registry
	s.provider = NewProvider(viper.GetString(store.TypeViperKey), memory.NewChainRegistry(), 200*time.Millisecond)
}

func (s *ProviderTestSuite) TestProvide() {
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
