package prometheus

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/standard"
)

const httpMetricsNamespace = "http"

type HTTP struct {
	*dynamicMulti
	*standard.HTTP
}

func NewHTTP(cfg *Config) *HTTP {
	multi := newDynamicMulti()

	reqsCounter := newCounter(
		multi,
		&prometheus.CounterOpts{
			Namespace: httpMetricsNamespace,
			Name:      "requests_total",
			Help:      "Total count of requests processed on a HTTP service",
		},
		[]string{"tenant_id", "entrypoint", "protocol", "service", "method", "code"},
	)
	multi.describers = append(multi.describers, reqsCounter.Describe)

	tlsReqsCounter := newCounter(
		multi,
		&prometheus.CounterOpts{
			Namespace: httpMetricsNamespace,
			Name:      "requests_tls_total",
			Help:      "Total count of TLS requests processed on a HTTP service",
		},
		[]string{"tenant_id", "entrypoint", "service", "tls_version", "tls_cipher"},
	)
	multi.describers = append(multi.describers, tlsReqsCounter.Describe)

	reqsLatencyHistogram := newHistogram(
		multi,
		&prometheus.HistogramOpts{
			Namespace: httpMetricsNamespace,
			Name:      "requests_latency_seconds",
			Help:      "Histogram of service's response latency (seconds)",
			Buckets:   cfg.HTTP.Buckets,
		},
		[]string{"tenant_id", "entrypoint", "protocol", "service", "method", "code"},
	)
	multi.describers = append(multi.describers, reqsLatencyHistogram.Describe)

	openConnsGauge := newGauge(
		multi,
		&prometheus.GaugeOpts{
			Namespace: httpMetricsNamespace,
			Name:      "open_connections",
			Help:      "Current count of open connections on a service",
		},
		[]string{"tenant_id", "entrypoint", "protocol", "service", "method"},
	)
	multi.describers = append(multi.describers, openConnsGauge.Describe)

	retriesCounter := newCounter(
		multi,
		&prometheus.CounterOpts{
			Namespace: httpMetricsNamespace,
			Name:      "retries_total",
			Help:      "Total count of request retries on a service.",
		},
		[]string{"tenant_id", "entrypoint", "service"},
	)
	multi.describers = append(multi.describers, retriesCounter.Describe)

	serverUpGauge := newGauge(
		multi,
		&prometheus.GaugeOpts{
			Namespace: httpMetricsNamespace,
			Name:      "server_up",
			Help:      "Current server status (0=DOWN, 1=UP)",
		},
		[]string{"tenant_id", "entrypoint", "service", "url"},
	)
	multi.describers = append(multi.describers, serverUpGauge.Describe)

	return &HTTP{
		dynamicMulti: multi,
		HTTP: standard.NewHTTP(
			reqsCounter,
			tlsReqsCounter,
			reqsLatencyHistogram,
			openConnsGauge,
			retriesCounter,
			serverUpGauge,
			func(*dynamic.Configuration) error { return nil },
		),
	}
}

func (http *HTTP) Switch(cfg *dynamic.Configuration) error {
	return http.dynamicMulti.Switch(cfg)
}

type dynamicMulti struct {
	describers []func(ch chan<- *prometheus.Desc)

	mux           *sync.Mutex
	currentCfg    *dynamicConfig
	currentValues map[string]*collector
	values        chan *collector
}

func newDynamicMulti() *dynamicMulti {
	cllctr := &dynamicMulti{
		mux:           &sync.Mutex{},
		currentCfg:    newDynamicConfig(nil),
		currentValues: make(map[string]*collector),
		values:        make(chan *collector, 100),
	}

	go cllctr.listenValues()

	return cllctr
}

// Describe implements prometheus.Collector and simply calls
// the registered describer functions.
func (c *dynamicMulti) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.describers {
		desc(ch)
	}
}

// Collect current collectors and then removes collectors that have
// removed of the dynamic configuration
func (c *dynamicMulti) Collect(ch chan<- prometheus.Metric) {
	c.mux.Lock()
	defer c.mux.Unlock()

	var outdatedKeys []string
	for key, value := range c.currentValues {
		value.Collect(ch)

		if c.isOutdated(value) {
			outdatedKeys = append(outdatedKeys, key)
		}
	}

	for _, key := range outdatedKeys {
		defer collectorPool.Put(c.currentValues[key])

		c.currentValues[key].delete()
		delete(c.currentValues, key)
	}
}

func (c *dynamicMulti) Switch(cfg *dynamic.Configuration) error {
	c.mux.Lock()
	c.currentCfg = newDynamicConfig(cfg)
	c.mux.Unlock()
	return nil
}

func (c *dynamicMulti) listenValues() {
	for value := range c.values {
		c.listenValue(value)
	}
}

func (c *dynamicMulti) listenValue(value *collector) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if oldValue, ok := c.currentValues[value.id]; ok {
		defer collectorPool.Put(oldValue)
	}

	c.currentValues[value.id] = value
}

// isOutdated checks whether the passed collector has labels that mark
// it as belonging to an outdated dynamic configuration
func (c *dynamicMulti) isOutdated(value *collector) bool {
	labels := value.labels

	if serviceName, ok := labels["service"]; ok {
		if !c.currentCfg.hasService(serviceName) {
			return true
		}

		if url, ok := labels["url"]; ok && !c.currentCfg.hasServerURL(serviceName, url) {
			return true
		}
	}

	return false
}

type vec interface {
	prometheus.Collector
	// With(prometheus.Labels) prometheus.Collector
	Delete(prometheus.Labels) bool
}

type collector struct {
	id            string
	name          string
	labels        prometheus.Labels
	promCollector prometheus.Collector
	delete        func() bool
	values        chan<- *collector
}

// collectors are pooled.
var collectorPool = sync.Pool{
	New: func() interface{} { return &collector{} },
}

func newCollector(id, name string, labels prometheus.Labels, promCollector prometheus.Collector, deleteMetric func() bool, values chan<- *collector) *collector {
	cllctr := collectorPool.Get().(*collector)
	cllctr.id = id
	cllctr.name = name
	if labels == nil {
		labels = prometheus.Labels{}
	}
	cllctr.labels = labels
	cllctr.promCollector = promCollector
	cllctr.delete = deleteMetric
	cllctr.values = values
	return cllctr
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.promCollector.Collect(ch)
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	c.promCollector.Describe(ch)
}

func (c *collector) withLabelValues(labelValues ...string) *collector {
	vector, ok := c.promCollector.(vec)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a vector)", c.promCollector))
	}

	labels := prometheus.Labels{}
	for k, v := range c.labels {
		labels[k] = v
	}

	for k, v := range ToLabels(labelValues...) {
		labels[k] = v
	}

	return newCollector(
		CollectorID(c.name, labels),
		c.name,
		labels,
		c.promCollector,
		func() bool { return vector.Delete(labels) },
		c.values,
	)
}

func (c *collector) with() interface{} {
	counterVec, ok := c.promCollector.(*prometheus.CounterVec)
	if ok {
		return counterVec.With(c.labels)
	}

	gaugeVec, ok := c.promCollector.(*prometheus.GaugeVec)
	if ok {
		return gaugeVec.With(c.labels)
	}

	histVec, ok := c.promCollector.(*prometheus.HistogramVec)
	if ok {
		return histVec.With(c.labels)
	}

	panic(fmt.Sprintf("invalid prometheus collector %T (not a vector)", c.promCollector))
}

func (c *collector) Add(delta float64) {
	cllctr := c.with()

	counter, ok := cllctr.(prometheus.Counter)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a %T)", cllctr, counter))
	}
	counter.Add(delta)

	c.newValue(counter)
}

func (c *collector) Set(value float64) {
	cllctr := c.with()

	gauge, ok := cllctr.(prometheus.Gauge)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a %T)", cllctr, gauge))
	}
	gauge.Set(value)

	c.newValue(gauge)
}

func (c *collector) Observe(value float64) {
	cllctr := c.with()

	observer, ok := cllctr.(prometheus.Histogram)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a %T)", cllctr, observer))
	}
	observer.Observe(value)

	c.newValue(observer)
}

func (c *collector) newValue(value prometheus.Collector) {
	c.values <- newCollector(
		c.id,
		"",
		c.labels,
		value,
		c.delete,
		nil,
	)
}

func CollectorID(name string, labels prometheus.Labels) string {
	var labelNamesValues []string
	for name, value := range labels {
		labelNamesValues = append(labelNamesValues, name, value)
	}
	sort.Strings(labelNamesValues)
	return name + ":" + strings.Join(labelNamesValues, "|")
}

type counter struct {
	*collector
}

func newCounter(multi *dynamicMulti, opts *prometheus.CounterOpts, labelNames []string) *counter {
	c := &counter{
		newCollector(
			"", opts.Name,
			nil,
			prometheus.NewCounterVec(*opts, labelNames),
			nil,
			multi.values,
		),
	}

	if len(labelNames) == 0 {
		c.Add(0)
	}

	return c
}

func (c *counter) With(labelValues ...string) kitmetrics.Counter {
	return &counter{c.collector.withLabelValues(labelValues...)}
}

type gauge struct {
	*collector
}

func newGauge(multi *dynamicMulti, opts *prometheus.GaugeOpts, labelNames []string) *gauge {
	g := &gauge{
		newCollector(
			"", opts.Name,
			nil,
			prometheus.NewGaugeVec(*opts, labelNames),
			nil,
			multi.values,
		),
	}

	if len(labelNames) == 0 {
		g.Set(0)
	}

	return g
}

func (g *gauge) With(labelValues ...string) kitmetrics.Gauge {
	return &gauge{g.collector.withLabelValues(labelValues...)}
}

type histogram struct {
	*collector
}

func newHistogram(multi *dynamicMulti, opts *prometheus.HistogramOpts, labelNames []string) *histogram {
	return &histogram{
		newCollector(
			"", opts.Name,
			nil,
			prometheus.NewHistogramVec(*opts, labelNames),
			nil,
			multi.values,
		),
	}
}

func (h *histogram) With(labelValues ...string) kitmetrics.Histogram {
	return &histogram{h.collector.withLabelValues(labelValues...)}
}

// dynamicConfig holds the current set of routers, and services
// corresponding to the current state of an HTTP router
// It provides a performant way to check whether the collected metrics belong to the
// current configuration or to an outdated one.
type dynamicConfig struct {
	routers  map[string]bool
	services map[string]map[string]bool
}

func newDynamicConfig(conf *dynamic.Configuration) *dynamicConfig {
	cfg := &dynamicConfig{
		routers:  make(map[string]bool),
		services: make(map[string]map[string]bool),
	}

	if conf == nil || conf.HTTP == nil {
		return cfg
	}

	for rtName := range conf.HTTP.Routers {
		cfg.routers[rtName] = true
	}

	for serviceName, service := range conf.HTTP.Services {
		cfg.services[serviceName] = make(map[string]bool)
		if service.ReverseProxy != nil {
			for _, server := range service.ReverseProxy.LoadBalancer.Servers {
				cfg.services[serviceName][server.URL] = true
			}
		}
	}

	return cfg
}

func (cfg *dynamicConfig) hasService(serviceName string) bool {
	_, ok := cfg.services[serviceName]
	return ok
}

func (cfg *dynamicConfig) hasServerURL(serviceName, serverURL string) bool {
	if service, hasService := cfg.services[serviceName]; hasService {
		_, ok := service[serverURL]
		return ok
	}
	return false
}
