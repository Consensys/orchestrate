package standard

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

type HTTP struct {
	reqsCounter         kitmetrics.Counter
	reqsTLSCounter      kitmetrics.Counter
	reqLatencyHistogram kitmetrics.Histogram
	openConnsGauge      kitmetrics.Gauge
	retriesCounter      kitmetrics.Counter
	serverUpGauge       kitmetrics.Gauge
	switchCfg           func(*dynamic.Configuration) error
}

func NewHTTP(
	reqsCounter, reqsTLSCounter kitmetrics.Counter,
	reqLatencyHistogram kitmetrics.Histogram,
	openConnsGauge kitmetrics.Gauge,
	retriesCounter kitmetrics.Counter,
	serverUpGauge kitmetrics.Gauge,
	switchCfg func(*dynamic.Configuration) error,
) *HTTP {
	return &HTTP{
		reqsCounter:         reqsCounter,
		reqsTLSCounter:      reqsTLSCounter,
		reqLatencyHistogram: reqLatencyHistogram,
		openConnsGauge:      openConnsGauge,
		retriesCounter:      retriesCounter,
		serverUpGauge:       serverUpGauge,
		switchCfg:           switchCfg,
	}
}

func (r *HTTP) RequestsCounter() kitmetrics.Counter {
	return r.reqsCounter
}

func (r *HTTP) TLSRequestsCounter() kitmetrics.Counter {
	return r.reqsTLSCounter
}

func (r *HTTP) RequestsLatencyHistogram() kitmetrics.Histogram {
	return r.reqLatencyHistogram
}

func (r *HTTP) OpenConnsGauge() kitmetrics.Gauge {
	return r.openConnsGauge
}

func (r *HTTP) RetriesCounter() kitmetrics.Counter {
	return r.retriesCounter
}

func (r *HTTP) ServerUpGauge() kitmetrics.Gauge {
	return r.serverUpGauge
}

func (r *HTTP) Switch(cfg *dynamic.Configuration) error {
	return r.switchCfg(cfg)
}
