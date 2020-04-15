// +build unit

package app

import (
	"context"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	mockprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/static"
	gogrpc "google.golang.org/grpc"
)

func TestApp(t *testing.T) {
	httpEps := http.NewEntryPoints(
		map[string]*traefikstatic.EntryPoint{
			"test-ep": {
				Address: "127.0.0.1:1",
				Transport: &traefikstatic.EntryPointsTransport{
					RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
					LifeCycle:          &traefikstatic.LifeCycle{},
				},
			},
		},
		static.NewBuilder(make(map[string]*router.Router)),
	)

	grpcEp := grpc.NewEntryPoint(
		&traefikstatic.EntryPoint{
			Address: "127.0.0.1:2",
			Transport: &traefikstatic.EntryPointsTransport{
				RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
				LifeCycle:          &traefikstatic.LifeCycle{},
			},
		},
		gogrpc.NewServer(),
	)

	w := configwatcher.New(
		&configwatcher.Config{ProvidersThrottleDuration: time.Second},
		mockprovider.New(),
		func(map[string]interface{}) interface{} { return nil },
		[]func(context.Context, interface{}) error{},
	)

	app := New(w, httpEps, grpcEp)

	err := app.Start(context.Background())
	require.NoError(t, err, "App should have started properly")

	// Wait for application to properly start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = app.Stop(ctx)
	assert.NoError(t, err, "App should have stop properly")
}
