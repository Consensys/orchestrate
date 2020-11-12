// +build unit

package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/mock"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp/metrics/mock"
	"google.golang.org/grpc"
)

func TestEntryPoint(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	cfg := &traefikstatic.EntryPoint{
		Address: "127.0.0.1:0",
		Transport: &traefikstatic.EntryPointsTransport{
			RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
			LifeCycle:          &traefikstatic.LifeCycle{},
		},
	}

	builder := mockserver.NewMockBuilder(ctrlr)

	reg := mock.NewMockTPCMetrics(ctrlr)
	acceptConnsCounter := mock2.NewMockCounter(ctrlr)
	closedConnsCounter := mock2.NewMockCounter(ctrlr)
	openConnsGauce := mock2.NewMockGauge(ctrlr)
	connsLatencyHisto := mock2.NewMockHistogram(ctrlr)
	ep := NewEntryPoint("", cfg, builder, reg)

	reg.EXPECT().AcceptedConnsCounter().Return(acceptConnsCounter)
	acceptConnsCounter.EXPECT().With(gomock.Any()).Return(acceptConnsCounter)
	acceptConnsCounter.EXPECT().Add(gomock.Any())

	reg.EXPECT().OpenConnsGauge().Return(openConnsGauce)
	openConnsGauce.EXPECT().With(gomock.Any()).Return(openConnsGauce)
	openConnsGauce.EXPECT().Add(gomock.Any())

	reg.EXPECT().ClosedConnsCounter().Return(closedConnsCounter)
	closedConnsCounter.EXPECT().With(gomock.Any()).Return(closedConnsCounter)

	reg.EXPECT().ConnsLatencyHistogram().Return(connsLatencyHisto)
	connsLatencyHisto.EXPECT().With(gomock.Any()).Return(connsLatencyHisto)
	
	builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(grpc.NewServer(), nil)
	_ = ep.BuildServer(context.Background(), nil)

	done := make(chan struct{})
	go func() {
		_ = ep.ListenAndServe(context.Background())
		close(done)
	}()

	// Wait a few millisecond for server to start
	time.Sleep(500 * time.Millisecond)

	_, err := net.Dial("tcp", ep.Addr())
	require.NoError(t, err, "Dial should not error")

	_ = ep.Shutdown(context.Background())
	<-done
	_ = ep.Close()
}
