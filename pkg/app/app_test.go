package app

import (
	"context"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher"
	mockwatcher "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/mock"
	mockhttp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/router/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/registry"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
)

func newTestConfig() *Config {
	return &Config{
		HTTP: &HTTP{
			EntryPoints: map[string]*traefikstatic.EntryPoint{
				"test-ep": {
					Address: "127.0.0.1:1",
					Transport: &traefikstatic.EntryPointsTransport{
						RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
						LifeCycle:          &traefikstatic.LifeCycle{},
					},
				},
			},
		},
		Watcher: &configwatcher.Config{
			ProvidersThrottleDuration: time.Millisecond,
		},
		Metrics: registry.NewConfig(viper.GetViper()),
	}
}

func TestApp(t *testing.T) {
	ctx := context.Background()
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	httpBuilder := mockhttp.NewMockBuilder(ctrlr)

	watcher := mockwatcher.NewMockWatcher(ctrlr)

	reg := mock.NewMockRegistry(ctrlr)
	app := newApp(newTestConfig(), httpBuilder, watcher, reg, log.NewLogger())

	reg.EXPECT().Add(gomock.AssignableToTypeOf(tcpmetrics.NewTCPMetrics(nil)))
	watcher.EXPECT().AddListener(gomock.Any()).Times(2)
	watcher.EXPECT().Run(gomock.Any())

	err := app.Start(ctx)
	require.NoError(t, err, "App should have started properly")

	// Wait for application to properly start
	time.Sleep(100 * time.Millisecond)

	cctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	err = app.Stop(cctx)
	assert.NoError(t, err, "App should have stop properly")

	watcher.EXPECT().Close()
	err = app.Close()
	assert.NoError(t, err, "App should have closed properly")
}
