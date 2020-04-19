package standard

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)

type TCP struct {
	acceptedConnsCounter  kitmetrics.Counter
	closedConnsCounter    kitmetrics.Counter
	connsLatencyHistogram kitmetrics.Histogram
	openConnsGauge        kitmetrics.Gauge
}

func NewTCP(
	acceptedConnsCounter, closedConnsCounter kitmetrics.Counter,
	connsLatencyHistogram kitmetrics.Histogram,
	openConnsGauge kitmetrics.Gauge,
) *TCP {
	return &TCP{
		acceptedConnsCounter:  acceptedConnsCounter,
		closedConnsCounter:    closedConnsCounter,
		connsLatencyHistogram: connsLatencyHistogram,
		openConnsGauge:        openConnsGauge,
	}
}

func (r *TCP) AcceptedConnsCounter() kitmetrics.Counter {
	return r.acceptedConnsCounter
}

func (r *TCP) ClosedConnsCounter() kitmetrics.Counter {
	return r.closedConnsCounter
}

func (r *TCP) ConnsLatencyHistogram() kitmetrics.Histogram {
	return r.connsLatencyHistogram
}

func (r *TCP) OpenConnsGauge() kitmetrics.Gauge {
	return r.openConnsGauge
}
