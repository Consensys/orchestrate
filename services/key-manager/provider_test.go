// +build unit

package keymanager

import (
	"math"
	"reflect"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestNewInternalConfig(t *testing.T) {
	expectedCfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"key-manager": {
					Router: &traefikdynamic.Router{
						EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
						Service:     "key-manager",
						Priority:    math.MaxInt32,
						Rule:        "PathPrefix(`/ethereum`) || PathPrefix(`/zk-snarks`)",
						Middlewares: []string{"base@logger-base"},
					},
				},
			},
			Middlewares: make(map[string]*dynamic.Middleware),
			Services: map[string]*dynamic.Service{
				"key-manager": {
					KeyManager: &dynamic.KeyManager{},
				},
			},
		},
	}

	assert.True(t, reflect.DeepEqual(NewInternalConfig(), expectedCfg), "Configuration should be identical")
}
