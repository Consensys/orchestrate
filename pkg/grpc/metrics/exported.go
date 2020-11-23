package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics"
)

//go:generate mockgen -source=exported.go -destination=mock/mock.go -package=mock

const ModuleName = "grpc"

type GRPCMetrics interface {
	StartedCounter() kitmetrics.Counter
	HandledCounter() kitmetrics.Counter
	StreamMsgReceivedCounter() kitmetrics.Counter
	StreamMsgSentCounter() kitmetrics.Counter
	HandledDurationHistogram() kitmetrics.Histogram
	pkgmetrics.Prometheus
}
