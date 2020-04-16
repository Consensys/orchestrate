// +build unit

package configwatcher

import (
	"math"
	"reflect"
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

func TestNewInternalConfig(t *testing.T) {
	expectedCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"transactions": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Service:     "transactions",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/transactions`)",
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
				"healthcheck": &dynamic.Router{
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"metrics"},
						Service:     "healthcheck",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/`)",
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				"strip-api": &dynamic.Middleware{
					Middleware: &traefikdynamic.Middleware{
						StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
							Regex: []string{"/api"},
						},
					},
				},
				"auth": &dynamic.Middleware{
					Auth: &dynamic.Auth{},
				},
				"base-accesslog": &dynamic.Middleware{
					AccessLog: &dynamic.AccessLog{
						Format: "json",
					},
				},
			},
			Services: map[string]*dynamic.Service{
				"transactions": &dynamic.Service{
					Transactions: &dynamic.Transactions{},
				},
				"dashboard": &dynamic.Service{
					Dashboard: &dynamic.Dashboard{},
				},
				"healthcheck": &dynamic.Service{
					HealthCheck: &dynamic.HealthCheck{},
				},
				"swagger": &dynamic.Service{
					Swagger: &dynamic.Swagger{
						SpecsFile: "./public/swagger-specs/types/transaction-scheduler/swagger.json",
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

	watcherCfg := &configwatcher.Config{}
	cfg := NewConfig(staticCfg, watcherCfg)
	assert.True(t, reflect.DeepEqual(
		cfg.DynamicCfg(), expectedCfg), "Configuration should be identical")
}
