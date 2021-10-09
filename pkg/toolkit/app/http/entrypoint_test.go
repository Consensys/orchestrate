// +build unit
// +build !race

package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/mock"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/router"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/router/static"
	mock3 "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/mock"
	mock2 "github.com/consensys/orchestrate/pkg/toolkit/tcp/metrics/mock"
	"github.com/consensys/orchestrate/pkg/toolkit/tls/generate"
	traefikstatic "github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type okHandler struct {
	next http.Handler
}

func (h *okHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.next.ServeHTTP(rw, req)
	rw.WriteHeader(http.StatusOK)
}

func TestEntryPoints(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	httpHandler := mock.NewMockHandler(ctrlr)
	httpsHandler := mock.NewMockHandler(ctrlr)

	// Prepare TLS configuration
	cert, err := generate.DefaultCertificate()
	require.NoError(t, err)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
	}

	// Create Router Configuration
	confs := map[string]*router.Router{
		"test-ep": {
			HTTP:      &okHandler{httpHandler},
			HTTPS:     &okHandler{httpsHandler},
			TLSConfig: tlsConfig,
		},
	}

	reg := mock2.NewMockTPCMetrics(ctrlr)
	acceptConnsCounter := mock3.NewMockCounter(ctrlr)
	closedConnsCounter := mock3.NewMockCounter(ctrlr)
	openConnsGauce := mock3.NewMockGauge(ctrlr)
	connsLatencyHisto := mock3.NewMockHistogram(ctrlr)

	reg.EXPECT().AcceptedConnsCounter().Times(2).Return(acceptConnsCounter)
	acceptConnsCounter.EXPECT().With(gomock.Any()).Times(2).Return(acceptConnsCounter)
	acceptConnsCounter.EXPECT().Add(gomock.Any()).Times(2)

	reg.EXPECT().OpenConnsGauge().Times(2).Return(openConnsGauce)
	openConnsGauce.EXPECT().With(gomock.Any()).Times(2).Return(openConnsGauce)
	openConnsGauce.EXPECT().Add(gomock.Any()).AnyTimes()

	reg.EXPECT().ClosedConnsCounter().Times(2).Return(closedConnsCounter)
	closedConnsCounter.EXPECT().With(gomock.Any()).Times(2).Return(closedConnsCounter)
	closedConnsCounter.EXPECT().Add(gomock.Any()).AnyTimes()

	reg.EXPECT().ConnsLatencyHistogram().Times(2).Return(connsLatencyHisto)
	connsLatencyHisto.EXPECT().With(gomock.Any()).Times(2).Return(connsLatencyHisto)
	connsLatencyHisto.EXPECT().Observe(gomock.Any()).AnyTimes()

	eps := NewEntryPoints(
		map[string]*traefikstatic.EntryPoint{
			"test-ep": {
				Address: "127.0.0.1:0",
				Transport: &traefikstatic.EntryPointsTransport{
					RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
					LifeCycle:          &traefikstatic.LifeCycle{},
				},
			},
		},
		static.NewBuilder(confs),
		reg,
	)
	_ = eps.Switch(context.Background(), nil)

	done := make(chan struct{})
	go func() {
		_ = eps.ListenAndServe(context.Background())
		close(done)
	}()

	// Wait a few millisecond for server to start
	time.Sleep(50 * time.Millisecond)

	url := eps.Addresses()["test-ep"]
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Test calling HTTP
	httpHandler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
	resp, err := client.Get(fmt.Sprintf("http://%v", url))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should have correct status")

	// Test calling HTTPS
	httpsHandler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
	resp, err = client.Get(fmt.Sprintf("https://%v", url))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should have correct status")

	err = eps.Shutdown(context.Background())
	assert.NoError(t, err)

	err = eps.Close()
	assert.NoError(t, err)
	<-done
}

func TestEntryPointsError(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()
	httpHandler := mock.NewMockHandler(ctrlr)

	// Create Router Configuration
	confs := map[string]*router.Router{
		"test-ep1": {HTTP: &okHandler{httpHandler}},
		"test-ep2": {HTTP: &okHandler{httpHandler}},
	}

	// We try to open 2 entry points on the same IP
	eps := NewEntryPoints(
		map[string]*traefikstatic.EntryPoint{
			"test-ep1": {
				Address: "127.0.0.1:10",
				Transport: &traefikstatic.EntryPointsTransport{
					RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
					LifeCycle:          &traefikstatic.LifeCycle{},
				},
			},
			"test-ep2": {
				Address: "127.0.0.1:10",
				Transport: &traefikstatic.EntryPointsTransport{
					RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
					LifeCycle:          &traefikstatic.LifeCycle{},
				},
			},
		},
		static.NewBuilder(confs),
		mock2.NewMockTPCMetrics(ctrlr),
	)
	_ = eps.Switch(context.Background(), nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	errors := eps.ListenAndServe(ctx)
	select {
	case <-errors:
	case <-time.After(time.Second):
		t.Errorf("Entrypoints should have error")
	}

	err := eps.Shutdown(ctx)
	assert.NoError(t, err)

	err = eps.Close()
	assert.NoError(t, err)
}
