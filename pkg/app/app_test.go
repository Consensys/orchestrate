// +build unit

package app

import (
	"context"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	mockprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/mock"
	mockgrpc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	mockhttp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/mock"
	mockmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/mock"
	gogrpc "google.golang.org/grpc"
)

func TestApp(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	cfg := &Config{
		HTTP: &traefikstatic.Configuration{
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
		},
		Watcher: &configwatcher.Config{
			ProvidersThrottleDuration: time.Millisecond,
		},
	}

	grpcBuilder := mockgrpc.NewMockBuilder(ctrlr)
	grpcBuilder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(gogrpc.NewServer(), nil)

	httpBuilder := mockhttp.NewMockBuilder(ctrlr)
	httpBuilder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	metricsRegistry := mockmetrics.NewMockRegistry(ctrlr)
	tcpMetrics := mockmetrics.NewMockTCP(ctrlr)
	metricsRegistry.EXPECT().TCP().Return(tcpMetrics).Times(2)

	httpMetrics := mockmetrics.NewMockHTTP(ctrlr)
	metricsRegistry.EXPECT().HTTP().Return(httpMetrics)

	prvdr := mockprovider.New()

	app, err := New(cfg, prvdr, httpBuilder, grpcBuilder, metricsRegistry)
	require.NoError(t, err, "App should have been created properly")

	err = app.Start(context.Background())
	require.NoError(t, err, "App should have started properly")

	// Wait for application to properly start
	time.Sleep(100 * time.Millisecond)

	dynCfg := &dynamic.Configuration{}
	msg := dynamic.NewMessage("test-provider", dynCfg)

	httpMetrics.EXPECT().Switch(gomock.Any())
	_ = prvdr.ProvideMsg(context.Background(), msg)

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = app.Stop(ctx)
	assert.NoError(t, err, "App should have stop properly")
}
