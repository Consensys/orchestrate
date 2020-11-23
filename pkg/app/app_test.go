package app

import (
	"context"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher"
	mockwatcher "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/mock"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/config/static"
	mockgrpc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/server/mock"
	mockhttp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/router/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/registry"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
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
		Metrics: registry.NewConfig(viper.GetViper()),
	}
}

func TestApp(t *testing.T) {
	ctx := context.Background()
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	grpcBuilder := mockgrpc.NewMockBuilder(ctrlr)
	grpcBuilder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(gogrpc.NewServer(), nil)

	httpBuilder := mockhttp.NewMockBuilder(ctrlr)

	watcher := mockwatcher.NewMockWatcher(ctrlr)

	reg := mock.NewMockRegistry(ctrlr)
	app := newApp(newTestConfig(), httpBuilder, grpcBuilder, watcher, reg, logrus.New())

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
