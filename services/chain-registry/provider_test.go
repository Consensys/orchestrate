// +build unit

package chainregistry

import (
	"math"
	"reflect"
	"testing"
	"time"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

func TestNewInternalConfig(t *testing.T) {
	expectedCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"chains": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Service:     "chains",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/chains`) || PathPrefix(`/faucets`)",
						Middlewares: []string{"base-accesslog", "auth"},
					},
				},
				"swagger": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Service:     "swagger",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/swagger`)",
						Middlewares: []string{"base-accesslog"},
					},
				},
				"dashboard": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Service:     "dashboard",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/api`) || PathPrefix(`/dashboard`)",
						Middlewares: []string{"base-accesslog", "strip-api"},
					},
				},
				"healthcheck": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"metrics"},
						Service:     "healthcheck",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/`)",
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"strip-api": {
					Middleware: &traefikdynamic.Middleware{
						StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
							Regex: []string{"/api"},
						},
					},
				},
				"strip-path": {
					Middleware: &traefikdynamic.Middleware{
						StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
							Regex: []string{"/.+"},
						},
					},
				},
				"auth": {
					Auth: &dynamic.Auth{},
				},
				"ratelimit": {
					RateLimit: &dynamic.RateLimit{
						MaxDelay:     time.Second,
						DefaultDelay: 30 * time.Second,
						Cooldown:     30 * time.Second,
					},
				},
				"base-accesslog": {
					AccessLog: &dynamic.AccessLog{
						Format: "json",
					},
				},
				"chain-proxy-accesslog": {
					AccessLog: &dynamic.AccessLog{
						Filters: &dynamic.AccessLogFilters{
							StatusCodes: []string{"100-199", "400-428", "430-599"},
						},
						Format: "json",
					},
				},
			},
			Services: map[string]*dynamic.Service{
				"chains": {
					Chains: &dynamic.Chains{},
				},
				"dashboard": {
					Dashboard: &dynamic.Dashboard{},
				},
				"healthcheck": {
					HealthCheck: &dynamic.HealthCheck{},
				},
				"swagger": {
					Swagger: &dynamic.Swagger{
						SpecsFile: "./public/swagger-specs/types/chain-registry/swagger.json",
					},
				},
			},
		},
	}

	staticCfg := &traefikstatic.Configuration{
		Log: &traefiktypes.TraefikLog{
			Format: "json",
		},
	}
	assert.True(t, reflect.DeepEqual(NewInternalConfig(staticCfg), expectedCfg), "Configuration should be identical")
}

func TestNewChainsProxyConfig(t *testing.T) {
	testSet := []struct {
		chains      []*types.Chain
		expectedCfg func(*dynamic.Configuration) *dynamic.Configuration
	}{
		{
			[]*types.Chain{
				{
					UUID:     "0d60a85e-0b90-4482-a14c-108aea2557aa",
					Name:     "testChain",
					TenantID: "testTenantId",
					URLs: []string{
						"http://testURL1.com",
						"http://testURL2.com",
					},
				},
			},
			func(cfg *dynamic.Configuration) *dynamic.Configuration {
				cfg.HTTP.Routers["chain-0d60a85e-0b90-4482-a14c-108aea2557aa"] = &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Priority:    math.MaxInt32,
						Service:     "chain-0d60a85e-0b90-4482-a14c-108aea2557aa",
						Rule:        "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`)",
						Middlewares: []string{
							"chain-proxy-accesslog@internal",
							"auth@internal",
							"multitenancy-testTenantId",
							"strip-path@internal",
							"ratelimit@internal",
						},
					},
				}
				cfg.HTTP.Services["chain-0d60a85e-0b90-4482-a14c-108aea2557aa"] = &dynamic.Service{
					ReverseProxy: &dynamic.ReverseProxy{
						PassHostHeader: utils.Bool(false),
						LoadBalancer: &dynamic.LoadBalancer{
							Servers: []*dynamic.Server{
								{URL: "http://testURL1.com"},
								{URL: "http://testURL2.com"},
							},
						},
					},
				}

				cfg.HTTP.Middlewares["multitenancy-testTenantId"] = &dynamic.Middleware{
					MultiTenancy: &dynamic.MultiTenancy{
						Tenant: "testTenantId",
					},
				}

				return cfg
			},
		},
		{
			[]*types.Chain{
				{
					UUID:     "0d60a85e-0b90-4482-a14c-108aea2557aa",
					Name:     "testChain",
					TenantID: "testTenantId",
					URLs: []string{
						"http://testURL1.com",
						"http://testURL2.com",
					},
				},
				{
					UUID:     "39240e9f-ae09-4e95-9fd0-a712035c8ad7",
					Name:     "testChain2",
					TenantID: "testTenantId",
					URLs: []string{
						"http://testURL10.com",
						"http://testURL20.com",
					},
				},
			},
			func(cfg *dynamic.Configuration) *dynamic.Configuration {
				cfg.HTTP.Routers["chain-0d60a85e-0b90-4482-a14c-108aea2557aa"] = &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Priority:    math.MaxInt32,
						Service:     "chain-0d60a85e-0b90-4482-a14c-108aea2557aa",
						Rule:        "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`)",
						Middlewares: []string{
							"chain-proxy-accesslog@internal",
							"auth@internal",
							"multitenancy-testTenantId",
							"strip-path@internal",
							"ratelimit@internal",
						},
					},
				}
				cfg.HTTP.Routers["chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7"] = &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Priority:    math.MaxInt32,
						Service:     "chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7",
						Rule:        "Path(`/39240e9f-ae09-4e95-9fd0-a712035c8ad7`)",
						Middlewares: []string{
							"chain-proxy-accesslog@internal",
							"auth@internal",
							"multitenancy-testTenantId",
							"strip-path@internal",
							"ratelimit@internal",
						},
					},
				}

				cfg.HTTP.Services["chain-0d60a85e-0b90-4482-a14c-108aea2557aa"] = &dynamic.Service{
					ReverseProxy: &dynamic.ReverseProxy{
						PassHostHeader: utils.Bool(false),
						LoadBalancer: &dynamic.LoadBalancer{
							Servers: []*dynamic.Server{
								{URL: "http://testURL1.com"},
								{URL: "http://testURL2.com"},
							},
						},
					},
				}

				cfg.HTTP.Services["chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7"] = &dynamic.Service{
					ReverseProxy: &dynamic.ReverseProxy{
						PassHostHeader: utils.Bool(false),
						LoadBalancer: &dynamic.LoadBalancer{
							Servers: []*dynamic.Server{
								{URL: "http://testURL10.com"},
								{URL: "http://testURL20.com"},
							},
						},
					},
				}

				cfg.HTTP.Middlewares["multitenancy-testTenantId"] = &dynamic.Middleware{
					MultiTenancy: &dynamic.MultiTenancy{
						Tenant: "testTenantId",
					},
				}
				return cfg
			},
		},
	}

	for i, test := range testSet {
		cfg := NewChainsProxyConfig(test.chains)
		expectedCfg := test.expectedCfg(dynamic.NewConfig())
		t.Logf("%v", cfg)
		assert.Equal(t, expectedCfg, cfg, "Chain-registry - Store (%d/%d): expected %v but got %v", i+1, len(testSet), expectedCfg, cfg)
	}
}
