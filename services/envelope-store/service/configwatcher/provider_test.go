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
				"envelopes": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"http"},
						Service:     "envelopes",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/envelopes`)",
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
						Priority:    math.MaxInt32 - 1,
						Rule:        "PathPrefix(`/`)",
					},
				},
				"prometheus": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{"metrics"},
						Service:     "prometheus",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/metrics`)",
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
				"auth": {
					Auth: &dynamic.Auth{},
				},
				"base-accesslog": {
					AccessLog: &dynamic.AccessLog{
						Format: "json",
					},
				},
			},
			Services: map[string]*dynamic.Service{
				"envelopes": {
					Envelopes: &dynamic.Envelopes{},
				},
				"dashboard": {
					Dashboard: &dynamic.Dashboard{},
				},
				"healthcheck": {
					HealthCheck: &dynamic.HealthCheck{},
				},
				"prometheus": {
					Prometheus: &dynamic.Prometheus{},
				},
				"swagger": {
					Swagger: &dynamic.Swagger{
						SpecsFile: "./public/swagger-specs/types/envelope-store/store.swagger.json",
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
