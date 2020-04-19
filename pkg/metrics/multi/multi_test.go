package multi

import (
	"testing"

	kitmetrics "github.com/go-kit/kit/metrics"
	kitgeneric "github.com/go-kit/kit/metrics/generic"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/generic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/mock"
)

func testCounter(t *testing.T, multiCounter kitmetrics.Counter, genericCounter *kitgeneric.Counter, desc string) {
	multiCounter.Add(1)
	assert.Equal(t, float64(1), genericCounter.Value(), desc)
	genericCounter.ValueReset()
}

func testGauge(t *testing.T, multiGauge kitmetrics.Gauge, genericGauge *kitgeneric.Gauge, desc string) {
	multiGauge.Set(1)
	assert.Equal(t, float64(1), genericGauge.Value(), desc)
	multiGauge.Set(0)
}

func testHistogram(t *testing.T, multiHistogram kitmetrics.Histogram, genericHistogram *kitgeneric.Histogram, desc string) {
	multiHistogram.Observe(1)
	assert.Equal(t, float64(1), genericHistogram.Quantile(1), desc)
}

func TestMulti(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	reg := &Multi{}

	genericTCP := generic.NewTCP()
	mockTCP := mock.NewMockTCP(ctrlr)
	mockTCP.EXPECT().AcceptedConnsCounter().Return(genericTCP.AcceptedConnsCounter())
	mockTCP.EXPECT().ClosedConnsCounter().Return(genericTCP.ClosedConnsCounter())
	mockTCP.EXPECT().ConnsLatencyHistogram().Return(genericTCP.ConnsLatencyHistogram())
	mockTCP.EXPECT().OpenConnsGauge().Return(genericTCP.OpenConnsGauge())

	genericHTTP := generic.NewHTTP()
	mockHTTP := mock.NewMockHTTP(ctrlr)
	mockHTTP.EXPECT().RequestsCounter().Return(genericHTTP.RequestsCounter())
	mockHTTP.EXPECT().TLSRequestsCounter().Return(genericHTTP.TLSRequestsCounter())
	mockHTTP.EXPECT().RequestsLatencyHistogram().Return(genericHTTP.RequestsLatencyHistogram())
	mockHTTP.EXPECT().OpenConnsGauge().Return(genericHTTP.OpenConnsGauge())
	mockHTTP.EXPECT().RetriesCounter().Return(genericHTTP.RetriesCounter())
	mockHTTP.EXPECT().ServerUpGauge().Return(genericHTTP.ServerUpGauge())

	genericGRPCServer := generic.NewGRPCServer()
	mockGRPCServer := mock.NewMockGRPCServer(ctrlr)
	mockGRPCServer.EXPECT().StartedCounter().Return(genericGRPCServer.StartedCounter())
	mockGRPCServer.EXPECT().HandledCounter().Return(genericGRPCServer.HandledCounter())
	mockGRPCServer.EXPECT().StreamMsgReceivedCounter().Return(genericGRPCServer.StreamMsgReceivedCounter())
	mockGRPCServer.EXPECT().StreamMsgSentCounter().Return(genericGRPCServer.StreamMsgSentCounter())
	mockGRPCServer.EXPECT().HandledDurationHistogram().Return(genericGRPCServer.HandledDurationHistogram())

	reg.tcps = append(reg.tcps, mockTCP)
	reg.https = append(reg.https, mockHTTP)
	reg.grpcServers = append(reg.grpcServers, mockGRPCServer)

	tcp := reg.TCP()
	testCounter(t, tcp.AcceptedConnsCounter(), genericTCP.AcceptedConnsCounter().(*kitgeneric.Counter), "TCP/RequestsCounter")
	testCounter(t, tcp.ClosedConnsCounter(), genericTCP.ClosedConnsCounter().(*kitgeneric.Counter), "TCP/TLSRequestsCounter")
	testHistogram(t, tcp.ConnsLatencyHistogram(), genericTCP.ConnsLatencyHistogram().(*kitgeneric.Histogram), "TCP/RequestsLatencyHistogram")
	testGauge(t, tcp.OpenConnsGauge(), genericTCP.OpenConnsGauge().(*kitgeneric.Gauge), "TCP/OpenConnsGauge")

	http := reg.HTTP()
	testCounter(t, http.RequestsCounter(), genericHTTP.RequestsCounter().(*kitgeneric.Counter), "HTTP/RequestsCounter")
	testCounter(t, http.TLSRequestsCounter(), genericHTTP.TLSRequestsCounter().(*kitgeneric.Counter), "HTTP/TLSRequestsCounter")
	testHistogram(t, http.RequestsLatencyHistogram(), genericHTTP.RequestsLatencyHistogram().(*kitgeneric.Histogram), "HTTP/RequestsLatencyHistogram")
	testGauge(t, http.OpenConnsGauge(), genericHTTP.OpenConnsGauge().(*kitgeneric.Gauge), "HTTP/OpenConnsGauge")
	testCounter(t, http.RetriesCounter(), genericHTTP.RetriesCounter().(*kitgeneric.Counter), "HTTP/RetriesCounter")
	testGauge(t, http.ServerUpGauge(), genericHTTP.ServerUpGauge().(*kitgeneric.Gauge), "HTTP/ServerUpGauge")

	grpcServer := reg.GRPCServer()
	testCounter(t, grpcServer.StartedCounter(), genericGRPCServer.StartedCounter().(*kitgeneric.Counter), "GRPCServer/StartedCounter")
	testCounter(t, grpcServer.HandledCounter(), genericGRPCServer.HandledCounter().(*kitgeneric.Counter), "GRPCServer/HandledCounter")
	testCounter(t, grpcServer.StreamMsgReceivedCounter(), genericGRPCServer.StreamMsgReceivedCounter().(*kitgeneric.Counter), "GRPCServer/StreamMsgReceivedCounter")
	testCounter(t, grpcServer.StreamMsgSentCounter(), genericGRPCServer.StreamMsgSentCounter().(*kitgeneric.Counter), "GRPCServer/StreamMsgSentCounter")
	testHistogram(t, grpcServer.HandledDurationHistogram(), genericGRPCServer.HandledDurationHistogram().(*kitgeneric.Histogram), "GRPCServer/HandledDurationHistogram")
}
