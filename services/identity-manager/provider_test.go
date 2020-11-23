// +build unit

package identitymanager

import (
	"math"
	"reflect"
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
)

func TestNewInternalConfig(t *testing.T) {
	expectedCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"identities": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Service:     "identities",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/accounts`)",
						Middlewares: []string{"base@logger-base", "auth@multitenancy"},
					},
				},
			},
			Middlewares: make(map[string]*dynamic.Middleware),
			Services: map[string]*dynamic.Service{
				"identities": {
					Identity: &dynamic.Identity{},
				},
			},
		},
	}

	assert.True(t, reflect.DeepEqual(NewInternalConfig(), expectedCfg), "Configuration should be identical")
}
