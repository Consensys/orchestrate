package app

import (
	"context"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	mockwatcher "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/mock"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	mockgrpc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/mock"
	mockhttp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/mock"
	mockmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/mock"
	metrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
	gogrpc "google.golang.org/grpc"
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
		GRPC: &GRPC{
			EntryPoint: &traefikstatic.EntryPoint{
				Address: "127.0.0.1:2",
				Transport: &traefikstatic.EntryPointsTransport{
					RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
					LifeCycle:          &traefikstatic.LifeCycle{},
				},
			},
			Static: &grpcstatic.Configuration{},
		},
		Watcher: &configwatcher.Config{
			ProvidersThrottleDuration: time.Millisecond,
		},
		Metrics: &metrics.Config{},
	}
}

func TestApp(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	grpcBuilder := mockgrpc.NewMockBuilder(ctrlr)
	grpcBuilder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(gogrpc.NewServer(), nil)

	httpBuilder := mockhttp.NewMockBuilder(ctrlr)

	metricsRegistry := mockmetrics.NewMockRegistry(ctrlr)
	tcpMetrics := mockmetrics.NewMockTCP(ctrlr)

	watcher := mockwatcher.NewMockWatcher(ctrlr)

	app := newApp(newTestConfig(), httpBuilder, grpcBuilder, watcher, metricsRegistry, logrus.New())

	metricsRegistry.EXPECT().TCP().Return(tcpMetrics).Times(2)
	watcher.EXPECT().AddListener(gomock.Any()).Times(2)
	watcher.EXPECT().Run(gomock.Any())

	err := app.Start(context.Background())
	require.NoError(t, err, "App should have started properly")

	// Wait for application to properly start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = app.Stop(ctx)
	assert.NoError(t, err, "App should have stop properly")

	watcher.EXPECT().Close()
	err = app.Close()
	assert.NoError(t, err, "App should have closed properly")
}
