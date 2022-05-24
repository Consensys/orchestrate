// +build unit

package api

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/api/proxy"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestNewInternalConfig(t *testing.T) {
	expectedCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"api": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Service:     "api",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/transactions`) || PathPrefix(`/schedules`) || PathPrefix(`/jobs`) || PathPrefix(`/accounts`) || PathPrefix(`/faucets`) || PathPrefix(`/contracts`) || PathPrefix(`/chains`)",
						Middlewares: []string{"base@logger-base", "auth@multitenancy"},
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"strip-path": {
					Middleware: &traefikdynamic.Middleware{
						StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
							Regex: []string{`/proxy/chains/(?:tessera/)?(?:[a-zA-Z\d-]*)/?`},
						},
					},
				},
				"ratelimit": {
					RateLimit: &dynamic.RateLimit{
						MaxDelay:     time.Second,
						DefaultDelay: 30 * time.Second,
						Cooldown:     30 * time.Second,
					},
				},
				"chain-proxy-accesslog": {
					AccessLog: &dynamic.AccessLog{
						Filters: &dynamic.AccessLogFilters{
							StatusCodes: []string{"100-199", "400-428", "430-599"},
						},
					},
				},
			},
			Services: map[string]*dynamic.Service{
				"api": {
					API: &dynamic.API{},
				},
			},
		},
	}

	assert.True(t, reflect.DeepEqual(newInternalConfig(), expectedCfg), "Configuration should be identical")
}

func TestNewChainsProxyConfig(t *testing.T) {
	testSet := []struct {
		chains      []*entities.Chain
		expectedCfg func(*dynamic.Configuration) *dynamic.Configuration
	}{
		{
			[]*entities.Chain{
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
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Priority:    math.MaxInt32,
						Service:     "chain-0d60a85e-0b90-4482-a14c-108aea2557aa",
						Rule:        "Path(`/proxy/chains/0d60a85e-0b90-4482-a14c-108aea2557aa`)",
						Middlewares: []string{
							"chain-proxy-accesslog@internal",
							"auth@multitenancy",
							"auth-testTenantId@multitenancy",
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

				cfg.HTTP.Middlewares["auth-testTenantId@multitenancy"] = &dynamic.Middleware{
					MultiTenancy: &dynamic.MultiTenancy{
						Tenant: "testTenantId",
					},
				}

				return cfg
			},
		},
		{
			[]*entities.Chain{
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
					Name:     "testTesseraChain2",
					TenantID: "testTenantId",
					URLs: []string{
						"http://testURL10.com",
					},
					PrivateTxManager: &entities.PrivateTxManager{
						URL:  "http://testURL10.com/tessera",
						Type: entities.TesseraChainType,
					},
				},
			},
			func(cfg *dynamic.Configuration) *dynamic.Configuration {
				cfg.HTTP.Routers["chain-0d60a85e-0b90-4482-a14c-108aea2557aa"] = &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Priority:    math.MaxInt32,
						Service:     "chain-0d60a85e-0b90-4482-a14c-108aea2557aa",
						Rule:        "Path(`/proxy/chains/0d60a85e-0b90-4482-a14c-108aea2557aa`)",
						Middlewares: []string{
							"chain-proxy-accesslog@internal",
							"auth@multitenancy",
							"auth-testTenantId@multitenancy",
							"strip-path@internal",
							"ratelimit@internal",
						},
					},
				}
				cfg.HTTP.Routers["chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7"] = &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Priority:    math.MaxInt32,
						Service:     "chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7",
						Rule:        "Path(`/proxy/chains/39240e9f-ae09-4e95-9fd0-a712035c8ad7`)",
						Middlewares: []string{
							"chain-proxy-accesslog@internal",
							"auth@multitenancy",
							"auth-testTenantId@multitenancy",
							"strip-path@internal",
							"ratelimit@internal",
						},
					},
				}
				cfg.HTTP.Routers["tessera-chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7"] = &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Priority:    math.MaxInt32,
						Service:     "tessera-chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7",
						Rule:        "PathPrefix(`/proxy/chains/tessera/39240e9f-ae09-4e95-9fd0-a712035c8ad7`)",
						Middlewares: []string{
							"chain-proxy-accesslog@internal",
							"auth@multitenancy",
							"auth-testTenantId@multitenancy",
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
							},
						},
					},
				}

				cfg.HTTP.Services["tessera-chain-39240e9f-ae09-4e95-9fd0-a712035c8ad7"] = &dynamic.Service{
					ReverseProxy: &dynamic.ReverseProxy{
						PassHostHeader: utils.Bool(false),
						LoadBalancer: &dynamic.LoadBalancer{
							Servers: []*dynamic.Server{
								{URL: "http://testURL10.com/tessera"},
							},
						},
					},
				}

				cfg.HTTP.Middlewares["auth-testTenantId@multitenancy"] = &dynamic.Middleware{
					MultiTenancy: &dynamic.MultiTenancy{
						Tenant: "testTenantId",
					},
				}
				return cfg
			},
		},
	}

	for i, test := range testSet {
		cfg := proxy.NewProxyConfig(test.chains, nil, true)
		expectedCfg := test.expectedCfg(dynamic.NewConfig())
		assert.Equal(t, expectedCfg, cfg, "Chain-registry - Store (%d/%d): expected %v but got %v", i+1, len(testSet), expectedCfg, cfg)
	}
}
