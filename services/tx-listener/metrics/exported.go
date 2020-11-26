package metrics

import (
	"fmt"

	kitmetrics "github.com/go-kit/kit/metrics"
	pkgmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics"
)

//go:generate mockgen -source=exported.go -destination=mock/mock.go -package=mock

var ModuleName = fmt.Sprintf("%s_%s", pkgmetrics.Namespace, Subsystem)

type ListenerMetrics interface {
	BlockCounter() kitmetrics.Counter
	pkgmetrics.Prometheus
}
