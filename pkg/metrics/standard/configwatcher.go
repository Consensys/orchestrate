package standard

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type Watcher struct {
	reloadsCounter               kitmetrics.Counter
	reloadsFailureCounter        kitmetrics.Counter
	lastConfigReloadSuccessGauge kitmetrics.Gauge
	lastConfigReloadFailureGauge kitmetrics.Gauge
}

func (r *Watcher) ReloadsCounter() kitmetrics.Counter {
	return r.reloadsCounter
}

func (r *Watcher) ReloadsFailureCounter() kitmetrics.Counter {
	return r.reloadsFailureCounter
}

func (r *Watcher) LastReloadSuccessGauge() kitmetrics.Gauge {
	return r.lastConfigReloadSuccessGauge
}

func (r *Watcher) LastReloadFailureGauge() kitmetrics.Gauge {
	return r.lastConfigReloadFailureGauge
}
