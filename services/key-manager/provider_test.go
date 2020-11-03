// +build unit

package keymanager

import (
	"math"
	"reflect"
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

func TestNewInternalConfig(t *testing.T) {
	expectedCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"ethereum": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Service:     "ethereum",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/ethereum/accounts`)",
						Middlewares: []string{"base@logger-base"},
					},
				},
			},
			Middlewares: make(map[string]*dynamic.Middleware),
			Services: map[string]*dynamic.Service{
				"ethereum": {
					Signer: &dynamic.Signer{},
				},
			},
		},
	}

	assert.True(t, reflect.DeepEqual(NewInternalConfig(), expectedCfg), "Configuration should be identical")
}
