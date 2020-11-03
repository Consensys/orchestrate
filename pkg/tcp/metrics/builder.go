package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type metrics struct {
	acceptedConnsCounter  kitmetrics.Counter
	closedConnsCounter    kitmetrics.Counter
	connsLatencyHistogram kitmetrics.Histogram
	openConnsGauge        kitmetrics.Gauge
}

func buildMetrics(
	acceptedConnsCounter, closedConnsCounter kitmetrics.Counter,
	connsLatencyHistogram kitmetrics.Histogram,
	openConnsGauge kitmetrics.Gauge,
) *metrics {
	return &metrics{
		acceptedConnsCounter:  acceptedConnsCounter,
		closedConnsCounter:    closedConnsCounter,
		connsLatencyHistogram: connsLatencyHistogram,
		openConnsGauge:        openConnsGauge,
	}
}

func (r *metrics) AcceptedConnsCounter() kitmetrics.Counter {
	return r.acceptedConnsCounter
}

func (r *metrics) ClosedConnsCounter() kitmetrics.Counter {
	return r.closedConnsCounter
}

func (r *metrics) ConnsLatencyHistogram() kitmetrics.Histogram {
	return r.connsLatencyHistogram
}

func (r *metrics) OpenConnsGauge() kitmetrics.Gauge {
	return r.openConnsGauge
}
