// +build unit

package tcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock3 "github.com/ConsenSys/orchestrate/pkg/metrics/mock"
	"github.com/ConsenSys/orchestrate/pkg/tcp/metrics"
	mock2 "github.com/ConsenSys/orchestrate/pkg/tcp/metrics/mock"
)

func prepareEntryPoint(t *testing.T, handler Handler, reg metrics.TPCMetrics) *EntryPoint {
	ep := NewEntryPoint(
		"test",
		&traefikstatic.EntryPoint{
			Address: "127.0.0.1:0",
			Transport: &traefikstatic.EntryPointsTransport{
				RespondingTimeouts: &traefikstatic.RespondingTimeouts{},
				LifeCycle: &traefikstatic.LifeCycle{
					RequestAcceptGraceTimeout: 0,
					GraceTimeOut:              traefiktypes.Duration(5 * time.Second),
				},
			},
		},
		handler,
		reg,
	)

	return ep
}

func dial(ep *EntryPoint) (net.Conn, error) {
	return net.Dial("tcp", ep.Addr())
}

func firstConn(ep *EntryPoint) (net.Conn, error) {
	var conn net.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = dial(ep)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("entry point never started: %v", err)
	}

	return conn, nil
}

func testServing(t *testing.T, handler Handler, test func(t *testing.T, ep *EntryPoint, firstConn net.Conn)) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	reg := mock2.NewMockTPCMetrics(ctrlr)
	acceptConnsCounter := mock3.NewMockCounter(ctrlr)
	closedConnsCounter := mock3.NewMockCounter(ctrlr)
	openConnsGauce := mock3.NewMockGauge(ctrlr)
	connsLatencyHisto := mock3.NewMockHistogram(ctrlr)

	reg.EXPECT().AcceptedConnsCounter().AnyTimes().Return(acceptConnsCounter)
	acceptConnsCounter.EXPECT().With(gomock.Any()).AnyTimes().Return(acceptConnsCounter)
	acceptConnsCounter.EXPECT().Add(gomock.Any()).AnyTimes()
	
	reg.EXPECT().OpenConnsGauge().AnyTimes().Return(openConnsGauce)
	openConnsGauce.EXPECT().With(gomock.Any()).AnyTimes().Return(openConnsGauce)
	openConnsGauce.EXPECT().Add(gomock.Any()).AnyTimes()

	reg.EXPECT().ClosedConnsCounter().AnyTimes().Return(closedConnsCounter)
	closedConnsCounter.EXPECT().With(gomock.Any()).AnyTimes().Return(closedConnsCounter)
	closedConnsCounter.EXPECT().Add(gomock.Any()).AnyTimes()

	reg.EXPECT().ConnsLatencyHistogram().AnyTimes().Return(connsLatencyHisto)
	connsLatencyHisto.EXPECT().With(gomock.Any()).AnyTimes().Return(connsLatencyHisto)
	connsLatencyHisto.EXPECT().Observe(gomock.Any()).AnyTimes()

	ep := prepareEntryPoint(t, handler, reg)

	done := make(chan struct{})
	go func() {
		_ = ep.ListenAndServe(context.Background())
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)

	conn, err := firstConn(ep)
	require.NoError(t, err, "Entrypoint should have started properly")

	test(t, ep, conn)

	select {
	case <-done:
	default:
		_ = ep.Shutdown(context.Background())
		<-done
		ep.Close()
	}
}

func TestShutdown(t *testing.T) {
	handler := HandlerFunc(func(conn WriteCloser) {
		for {
			_, err := http.ReadRequest(bufio.NewReader(conn))
			if err == io.EOF || (err != nil && strings.HasSuffix(err.Error(), "use of closed network connection")) {
				return
			}
			require.NoError(t, err)

			resp := http.Response{StatusCode: http.StatusOK}
			err = resp.Write(conn)
			require.NoError(t, err)
		}
	})

	testServing(t, handler, func(t *testing.T, ep *EntryPoint, firstConn net.Conn) {
		// We need to do a write on a conn before the shutdown to make it "exist".
		// Because the connection indeed exists as far as TCP is concerned,
		// but since we only pass it along to the HTTP server after at least one byte is peaked,
		// the HTTP server (and hence its shutdown) does not know about the connection until that first byte peaking.
		request, _ := http.NewRequest(http.MethodGet, "", nil)
		err := request.Write(firstConn)
		require.NoError(t, err)

		go func() { _ = ep.Shutdown(context.Background()) }()

		// Make sure that new connections are not permitted anymore.
		// Note that this should be true not only after Shutdown has returned,
		// but technically also as early as the Shutdown has closed the listener,
		// i.e. during the shutdown and before the gracetime is over.
		var hasClosed bool
		var conn net.Conn
		for i := 0; i < 10; i++ {
			conn, err = dial(ep)
			if err == nil {
				conn.Close()
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if !strings.HasSuffix(err.Error(), "connection refused") && !strings.HasSuffix(err.Error(), "reset by peer") {
				t.Fatalf(`unexpected error: got %v, wanted "connection refused" or "reset by peer"`, err)
			}
			hasClosed = true
			break
		}
		require.True(t, hasClosed, "Entry point never closed")

		// And make sure that the connection we had opened before shutting things down is still operational
		resp, err := http.ReadResponse(bufio.NewReader(firstConn), request)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
