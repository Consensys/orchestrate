package internal

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/safe"
	"github.com/containous/traefik/v2/pkg/tls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProviderTestSuite struct {
	suite.Suite
	provider *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	s.provider = New()
}

func (s *ProviderTestSuite) TestProvide() {

	c := make(chan dynamic.Message)
	pool := safe.NewPool(context.Background())
	go func() {
		err := s.provider.Provide(c, pool)
		assert.NoError(s.T(), err, "Should Provide without error")
	}()
	config := <-c
	assert.Equal(s.T(), "internal", config.ProviderName, "Should get the correct providerName")

	pool.Stop()
}

func (s *ProviderTestSuite) TestInit() {
	assert.Nil(s.T(), s.provider.Init(), "should init without error")
}

func (s *ProviderTestSuite) TestCreateConfiguration() {
	c := s.provider.createConfiguration()

	cfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"api": {
					EntryPoints: []string{"orchestrate"},
					Service:     "api@internal",
					Priority:    math.MaxInt32 - 1,
					Rule:        "PathPrefix(`/{tenantID}`)",
					Middlewares: []string{"orchestrate-auth"},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"strip-path": &dynamic.Middleware{
					StripPrefixRegex: &dynamic.StripPrefixRegex{
						Regex: []string{"/.+"},
					},
				},
			},
			Services: map[string]*dynamic.Service{
				"api": {},
			},
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:  make(map[string]*dynamic.TCPRouter),
			Services: make(map[string]*dynamic.TCPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]tls.Store),
			Options: make(map[string]tls.Options),
		},
	}

	assert.True(s.T(), reflect.DeepEqual(c, cfg), "should be identical")
}

func (s *ProviderTestSuite) TestNew() {
	assert.NotNil(s.T(), s.provider, "should initialize provider")
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
