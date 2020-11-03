package metrics

import (
	"github.com/go-kit/kit/metrics/discard"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
)

const metricsNamespace = "grpc"

type grpcMetrics struct {
	prometheus.Collector
	*metrics
}

func NewGRPCMetrics(cfg *Config) GRPCMetrics {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}

	multi := pkgmetrics.NewMulti()

	startedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "started_total",
			Help:      "Total count of RPCs started on the server",
		},
		[]string{"tenant_id", "type", "service", "method"},
	)
	multi.Collectors = append(multi.Collectors, startedCounter)

	handledCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "handled_total",
			Help:      "Total count of RPCs completed on the server, regardless of success or failure.",
		},
		[]string{"tenant_id", "type", "service", "method", "code"},
	)
	multi.Collectors = append(multi.Collectors, handledCounter)

	streamMsgReceivedCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "msg_received_total",
			Help:      "Total count of RPC stream messages received on the server.",
		},
		[]string{"tenant_id", "type", "service", "method"},
	)
	multi.Collectors = append(multi.Collectors, streamMsgReceivedCounter)

	streamMsgSentCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "msg_sent_total",
			Help:      "Total count of RPC stream messages sent by the server.",
		},
		[]string{"tenant_id", "type", "service", "method"},
	)
	multi.Collectors = append(multi.Collectors, streamMsgSentCounter)

	HandledDurationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "handled_seconds",
			Help:      "Histogram of response latency (seconds) of RPCs handled by the server.",
			Buckets:   cfg.Buckets,
		},
		[]string{"tenant_id", "type", "service", "method", "code"},
	)
	multi.Collectors = append(multi.Collectors, HandledDurationHistogram)

	return &grpcMetrics{
		Collector: multi,
		metrics: buildMetrics(
			kitprometheus.NewCounter(startedCounter),
			kitprometheus.NewCounter(handledCounter),
			kitprometheus.NewCounter(streamMsgReceivedCounter),
			kitprometheus.NewCounter(streamMsgSentCounter),
			kitprometheus.NewHistogram(HandledDurationHistogram),
		),
	}
}

func NewGRPCNopMetrics(_ *Config) GRPCMetrics {
	return &grpcMetrics{
		Collector: pkgmetrics.NewMulti(),
		metrics: buildMetrics(
			discard.NewCounter(),
			discard.NewCounter(),
			discard.NewCounter(),
			discard.NewCounter(),
			discard.NewHistogram(),
		),
	}
}
