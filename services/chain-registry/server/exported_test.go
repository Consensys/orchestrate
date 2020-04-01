// +build unit

package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/containous/traefik/v2/pkg/provider/acme"
	"github.com/containous/traefik/v2/pkg/provider/aggregator"
	traefiktls "github.com/containous/traefik/v2/pkg/tls"

	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	SetGlobalStaticConfig(&static.Configuration{
		EntryPoints: static.EntryPoints{
			"test": {Address: "localhost:8080", ForwardedHeaders: &static.ForwardedHeaders{Insecure: true, TrustedIPs: []string{}}},
		},
		CertificatesResolvers: map[string]static.CertificateResolver{},
	})
	Init(context.Background())
	assert.NotNil(t, GlobalServer(), "Global server should have been set")

	var s *Server
	SetGlobalServer(s)
	assert.Nil(t, GlobalServer(), "Global should be reset to nil")
}

func TestInitACMEProvider(t *testing.T) {

	testSuite := []struct {
		c                  *static.Configuration
		providerAggregator *aggregator.ProviderAggregator
		tlsManager         *traefiktls.Manager
		expectedOutput     func() []*acme.Provider
	}{
		{
			&static.Configuration{
				CertificatesResolvers: map[string]static.CertificateResolver{},
			},
			&aggregator.ProviderAggregator{},
			traefiktls.NewManager(),
			func() []*acme.Provider {
				providers := make([]*acme.Provider, 0)
				return providers
			},
		},
		{
			&static.Configuration{
				CertificatesResolvers: map[string]static.CertificateResolver{
					"test": {},
				},
			},
			&aggregator.ProviderAggregator{},
			traefiktls.NewManager(),
			func() []*acme.Provider {
				providers := make([]*acme.Provider, 0)
				return providers
			},
		},
	}

	for _, test := range testSuite {
		output := initACMEProvider(test.c, test.providerAggregator, test.tlsManager)
		assert.True(t, compareProvider(output, test.expectedOutput()), "should be equal")
	}
}

// compareProvider is a partial substitute of reflect.DeepEqual to compare the output and the expected output of initACMEProvider as some unexported fields are set directly inside the function
func compareProvider(p1, p2 []*acme.Provider) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i := 0; i < len(p1); i++ {
		if !reflect.DeepEqual(p1[i].Configuration, p2[i].Configuration) {
			return false
		}
		if p1[i].ResolverName != p2[i].ResolverName {
			return false
		}
		if !reflect.DeepEqual(p1[i].ChallengeStore, p2[i].ChallengeStore) {
			return false
		}
	}
	return true
}
