package multi

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	kitmulti "github.com/go-kit/kit/metrics/multi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/standard"
)

type Multi struct {
	prometheus *prometheus.Prometheus

	tcps        []metrics.TCP
	https       []metrics.HTTP
	grpcServers []metrics.GRPCServer
}

func New(cfg *Config) *Multi {
	multi := &Multi{}

	if cfg.Prometheus != nil {
		multi.prometheus = prometheus.New(cfg.Prometheus)
		multi.tcps = append(multi.tcps, multi.prometheus.TCP())
		multi.https = append(multi.https, multi.prometheus.HTTP())
		multi.grpcServers = append(multi.grpcServers, multi.prometheus.GRPCServer())
	}

	return multi
}

func (multi *Multi) Prometheus() *prometheus.Prometheus {
	return multi.prometheus
}

func (multi *Multi) TCP() metrics.TCP {
	var acceptedConnsCounters []kitmetrics.Counter
	var closedConnsCounters []kitmetrics.Counter
	var connsLatencyHistogram []kitmetrics.Histogram
	var openConnsGauge []kitmetrics.Gauge

	for _, tcp := range multi.tcps {
		acceptedConnsCounters = append(acceptedConnsCounters, tcp.AcceptedConnsCounter())
		closedConnsCounters = append(closedConnsCounters, tcp.ClosedConnsCounter())
		connsLatencyHistogram = append(connsLatencyHistogram, tcp.ConnsLatencyHistogram())
		openConnsGauge = append(openConnsGauge, tcp.OpenConnsGauge())
	}

	return standard.NewTCP(
		kitmulti.NewCounter(acceptedConnsCounters...),
		kitmulti.NewCounter(closedConnsCounters...),
		kitmulti.NewHistogram(connsLatencyHistogram...),
		kitmulti.NewGauge(openConnsGauge...),
	)
}

func (multi *Multi) HTTP() metrics.HTTP {
	var acceptedConnsCounters []kitmetrics.Counter
	var closedConnsCounters []kitmetrics.Counter
	var connsLatencyHistogram []kitmetrics.Histogram
	var openConnsGauge []kitmetrics.Gauge
	var retriesCounter []kitmetrics.Counter
	var serverUpGauge []kitmetrics.Gauge
	var switches []func(*dynamic.Configuration) error

	for _, http := range multi.https {
		acceptedConnsCounters = append(acceptedConnsCounters, http.RequestsCounter())
		closedConnsCounters = append(closedConnsCounters, http.TLSRequestsCounter())
		connsLatencyHistogram = append(connsLatencyHistogram, http.RequestsLatencyHistogram())
		openConnsGauge = append(openConnsGauge, http.OpenConnsGauge())
		retriesCounter = append(retriesCounter, http.RetriesCounter())
		serverUpGauge = append(serverUpGauge, http.ServerUpGauge())
		switches = append(switches, http.Switch)
	}

	return standard.NewHTTP(
		kitmulti.NewCounter(acceptedConnsCounters...),
		kitmulti.NewCounter(closedConnsCounters...),
		kitmulti.NewHistogram(connsLatencyHistogram...),
		kitmulti.NewGauge(openConnsGauge...),
		kitmulti.NewCounter(retriesCounter...),
		kitmulti.NewGauge(serverUpGauge...),
		CombineSwitches(switches...),
	)
}

func (multi *Multi) GRPCServer() metrics.GRPCServer {
	var startedCounters []kitmetrics.Counter
	var handledCounters []kitmetrics.Counter
	var streamMsgReceivedCounters []kitmetrics.Counter
	var streamMsgSentCounters []kitmetrics.Counter
	var HandledDurationHistogram []kitmetrics.Histogram

	for _, grpcServer := range multi.grpcServers {
		startedCounters = append(startedCounters, grpcServer.StartedCounter())
		handledCounters = append(handledCounters, grpcServer.HandledCounter())
		streamMsgReceivedCounters = append(streamMsgReceivedCounters, grpcServer.StreamMsgReceivedCounter())
		streamMsgSentCounters = append(streamMsgSentCounters, grpcServer.StreamMsgSentCounter())
		HandledDurationHistogram = append(HandledDurationHistogram, grpcServer.HandledDurationHistogram())
	}

	return standard.NewGRPCServer(
		kitmulti.NewCounter(startedCounters...),
		kitmulti.NewCounter(handledCounters...),
		kitmulti.NewCounter(streamMsgReceivedCounters...),
		kitmulti.NewCounter(streamMsgSentCounters...),
		kitmulti.NewHistogram(HandledDurationHistogram...),
	)
}

func CombineSwitches(switches ...func(*dynamic.Configuration) error) func(*dynamic.Configuration) error {
	return func(cfg *dynamic.Configuration) error {
		for _, swtch := range switches {
			_ = swtch(cfg)
		}
		return nil
	}
}
