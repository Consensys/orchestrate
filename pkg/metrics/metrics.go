package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

//go:generate mockgen -source=metrics.go -destination=mock/mock.go -package=mock

type TCP interface {
	AcceptedConnsCounter() kitmetrics.Counter
	ClosedConnsCounter() kitmetrics.Counter
	ConnsLatencyHistogram() kitmetrics.Histogram
	OpenConnsGauge() kitmetrics.Gauge
}

type HTTP interface {
	RequestsCounter() kitmetrics.Counter
	TLSRequestsCounter() kitmetrics.Counter
	RequestsLatencyHistogram() kitmetrics.Histogram
	OpenConnsGauge() kitmetrics.Gauge
	RetriesCounter() kitmetrics.Counter
	ServerUpGauge() kitmetrics.Gauge
	Switch(cfg *dynamic.Configuration) error
}

type GRPCServer interface {
	StartedCounter() kitmetrics.Counter
	HandledCounter() kitmetrics.Counter
	StreamMsgReceivedCounter() kitmetrics.Counter
	StreamMsgSentCounter() kitmetrics.Counter
	HandledDurationHistogram() kitmetrics.Histogram
}

type Registry interface {
	TCP() TCP
	HTTP() HTTP
	GRPCServer() GRPCServer
}
