package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type metrics struct {
	reqsCounter         kitmetrics.Counter
	reqsTLSCounter      kitmetrics.Counter
	reqLatencyHistogram kitmetrics.Histogram
	openConnsGauge      kitmetrics.Gauge
	retriesCounter      kitmetrics.Counter
	serverUpGauge       kitmetrics.Gauge
}

func buildMetrics(
	reqsCounter, reqsTLSCounter kitmetrics.Counter,
	reqLatencyHistogram kitmetrics.Histogram,
	openConnsGauge kitmetrics.Gauge,
	retriesCounter kitmetrics.Counter,
	serverUpGauge kitmetrics.Gauge,
) *metrics {
	return &metrics{
		reqsCounter:         reqsCounter,
		reqsTLSCounter:      reqsTLSCounter,
		reqLatencyHistogram: reqLatencyHistogram,
		openConnsGauge:      openConnsGauge,
		retriesCounter:      retriesCounter,
		serverUpGauge:       serverUpGauge,
	}
}

func (r *metrics) RequestsCounter() kitmetrics.Counter {
	return r.reqsCounter
}

func (r *metrics) TLSRequestsCounter() kitmetrics.Counter {
	return r.reqsTLSCounter
}

func (r *metrics) RequestsLatencyHistogram() kitmetrics.Histogram {
	return r.reqLatencyHistogram
}

func (r *metrics) OpenConnsGauge() kitmetrics.Gauge {
	return r.openConnsGauge
}

func (r *metrics) RetriesCounter() kitmetrics.Counter {
	return r.retriesCounter
}

func (r *metrics) ServerUpGauge() kitmetrics.Gauge {
	return r.serverUpGauge
}
