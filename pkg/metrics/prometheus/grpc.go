package prometheus

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/standard"
)

const grpcServerMetricsNamespace = "grpc_server"

type GRPCServer struct {
	prometheus.Collector
	*standard.GRPCServer
}

func NewGRPCServer(cfg *Config) *GRPCServer {
	multi := NewMulti()

	startedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: grpcServerMetricsNamespace,
			Name:      "started_total",
			Help:      "Total count of RPCs started on the server",
		},
		[]string{"tenant_id", "type", "service", "method"},
	)
	multi.collectors = append(multi.collectors, startedCounter)

	handledCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: grpcServerMetricsNamespace,
			Name:      "handled_total",
			Help:      "Total count of RPCs completed on the server, regardless of success or failure.",
		},
		[]string{"tenant_id", "type", "service", "method", "code"},
	)
	multi.collectors = append(multi.collectors, handledCounter)

	streamMsgReceivedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: grpcServerMetricsNamespace,
			Name:      "msg_received_total",
			Help:      "Total count of RPC stream messages received on the server.",
		},
		[]string{"tenant_id", "type", "service", "method"},
	)
	multi.collectors = append(multi.collectors, streamMsgReceivedCounter)

	streamMsgSentCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: grpcServerMetricsNamespace,
			Name:      "msg_sent_total",
			Help:      "Total count of RPC stream messages sent by the server.",
		},
		[]string{"tenant_id", "type", "service", "method"},
	)
	multi.collectors = append(multi.collectors, streamMsgSentCounter)

	HandledDurationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: grpcServerMetricsNamespace,
			Name:      "handled_seconds",
			Help:      "Histogram of response latency (seconds) of RPCs handled by the server.",
			Buckets:   cfg.GRPC.Buckets,
		},
		[]string{"tenant_id", "type", "service", "method", "code"},
	)
	multi.collectors = append(multi.collectors, HandledDurationHistogram)

	return &GRPCServer{
		Collector: multi,
		GRPCServer: standard.NewGRPCServer(
			kitprometheus.NewCounter(startedCounter),
			kitprometheus.NewCounter(handledCounter),
			kitprometheus.NewCounter(streamMsgReceivedCounter),
			kitprometheus.NewCounter(streamMsgSentCounter),
			kitprometheus.NewHistogram(HandledDurationHistogram),
		),
	}
}
