// +build unit

package contractregistry

import (
	"math"
	"reflect"
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

func TestNewInternalConfig(t *testing.T) {
	expectedCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"contracts": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Service:     "contracts",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/contracts`)",
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
				"contracts": &dynamic.Service{
					Contracts: &dynamic.Contracts{},
				},
				"dashboard": &dynamic.Service{
					Dashboard: &dynamic.Dashboard{},
				},
				"healthcheck": &dynamic.Service{
					HealthCheck: &dynamic.HealthCheck{},
				},
				"swagger": &dynamic.Service{
					Swagger: &dynamic.Swagger{
						SpecsFile: "./public/swagger-specs/types/contract-registry/registry.swagger.json",
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
