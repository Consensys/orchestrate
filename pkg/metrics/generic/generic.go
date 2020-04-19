package generic

import (
	kitgeneric "github.com/go-kit/kit/metrics/generic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/standard"
)

const (
	TCPPrefix                   = "ep."
	TCPAcceptedConnsCounterName = TCPPrefix + "reqs.counter"
	TCPClosedConnsCounterName   = TCPPrefix + "tls-reqs.counter"
	TCPReqsLatencyHistogramName = TCPPrefix + "reqs-latency.histogram"
	TCPOpenConnsGaugeName       = TCPPrefix + "open-conns.gauge"
)

func NewTCP() metrics.TCP {
	return standard.NewTCP(
		kitgeneric.NewCounter(TCPAcceptedConnsCounterName),
		kitgeneric.NewCounter(TCPClosedConnsCounterName),
		kitgeneric.NewHistogram(TCPReqsLatencyHistogramName, 50),
		kitgeneric.NewGauge(TCPOpenConnsGaugeName),
	)
}

const (
	HTTPPrefix                   = "http."
	HTTPAcceptedConnsCounterName = HTTPPrefix + "reqs.counter"
	HTTPClosedConnsCounterName   = HTTPPrefix + "tls-reqs.counter"
	HTTPReqsLatencyHistogramName = HTTPPrefix + "reqs-latency.histogram"
	HTTPOpenConnsGaugeName       = HTTPPrefix + "open-conns.gauge"
	HTTPRetriesCounterName       = HTTPPrefix + "retries.counter"
	HTTPServerUpGaugeName        = HTTPPrefix + "server-up.gauge"
)

func NewHTTP() metrics.HTTP {
	return standard.NewHTTP(
		kitgeneric.NewCounter(HTTPAcceptedConnsCounterName),
		kitgeneric.NewCounter(HTTPClosedConnsCounterName),
		kitgeneric.NewHistogram(HTTPReqsLatencyHistogramName, 50),
		kitgeneric.NewGauge(HTTPOpenConnsGaugeName),
		kitgeneric.NewCounter(HTTPRetriesCounterName),
		kitgeneric.NewGauge(HTTPServerUpGaugeName),
		func(*dynamic.Configuration) error { return nil },
	)
}

const (
	GRPCServerPrefix                       = "gprc-server."
	GRPCServerStartedCounterName           = GRPCServerPrefix + "started.counter"
	GRPCServerHandledCounterName           = GRPCServerPrefix + "handled.counter"
	GRPCServerStreamMsgReceivedCounterName = GRPCServerPrefix + "msg-received.counter"
	GRPCServerStreamMsgSentCounterName     = GRPCServerPrefix + "msg-sent.counter"
	GRPCServerHandledLatencyHistogramName  = GRPCServerPrefix + "handled-latency.histogram"
)

func NewGRPCServer() metrics.GRPCServer {
	return standard.NewGRPCServer(
		kitgeneric.NewCounter(GRPCServerStartedCounterName),
		kitgeneric.NewCounter(GRPCServerHandledCounterName),
		kitgeneric.NewCounter(GRPCServerStreamMsgReceivedCounterName),
		kitgeneric.NewCounter(GRPCServerStreamMsgSentCounterName),
		kitgeneric.NewHistogram(GRPCServerHandledLatencyHistogramName, 50),
	)
}
