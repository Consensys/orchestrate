package envelopestore

import (
	"math"
	"reflect"
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
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
						Middlewares: []string{"base@logger-base", "auth@multitenancy"},
					},
				},
			},
			Middlewares: make(map[string]*dynamic.Middleware),
			Services: map[string]*dynamic.Service{
				"envelopes": {
					Envelopes: &dynamic.Envelopes{},
				},
			},
		},
	}

	assert.True(t, reflect.DeepEqual(NewInternalConfig(), expectedCfg), "Configuration should be identical")
}
