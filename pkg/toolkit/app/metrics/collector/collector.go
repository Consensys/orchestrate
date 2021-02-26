package collector

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type vec interface {
	prometheus.Collector
	// With(prometheus.Labels) prometheus.Collector
	Delete(prometheus.Labels) bool
}

type Collector struct {
	ID            string
	Name          string
	Labels        prometheus.Labels
	promCollector prometheus.Collector
	Delete        func() bool
	Values        chan<- *Collector
}

var CollectorPool = sync.Pool{
	New: func() interface{} { return &Collector{} },
}

func NewCollector(id, name string, labels prometheus.Labels, promCollector prometheus.Collector, deleteMetric func() bool, values chan<- *Collector) *Collector {
	cllctr := CollectorPool.Get().(*Collector)
	cllctr.ID = id
	cllctr.Name = name
	if labels == nil {
		labels = prometheus.Labels{}
	}
	cllctr.Labels = labels
	cllctr.promCollector = promCollector
	cllctr.Delete = deleteMetric
	cllctr.Values = values
	return cllctr
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.promCollector.Collect(ch)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.promCollector.Describe(ch)
}

func (c *Collector) withLabelValues(labelValues ...string) *Collector {
	vector, ok := c.promCollector.(vec)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a vector)", c.promCollector))
	}

	labels := prometheus.Labels{}
	for k, v := range c.Labels {
		labels[k] = v
	}

	for k, v := range toLabels(labelValues...) {
		labels[k] = v
	}

	return NewCollector(
		NewCollectorID(c.Name, labels),
		c.Name,
		labels,
		c.promCollector,
		func() bool { return vector.Delete(labels) },
		c.Values,
	)
}

func (c *Collector) with() interface{} {
	counterVec, ok := c.promCollector.(*prometheus.CounterVec)
	if ok {
		return counterVec.With(c.Labels)
	}

	gaugeVec, ok := c.promCollector.(*prometheus.GaugeVec)
	if ok {
		return gaugeVec.With(c.Labels)
	}

	histVec, ok := c.promCollector.(*prometheus.HistogramVec)
	if ok {
		return histVec.With(c.Labels)
	}

	panic(fmt.Sprintf("invalid prometheus collector %T (not a vector)", c.promCollector))
}

func (c *Collector) Add(delta float64) {
	cllctr := c.with()

	counter, ok := cllctr.(prometheus.Counter)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a %T)", cllctr, counter))
	}
	counter.Add(delta)

	c.newValue(counter)
}

func (c *Collector) Set(value float64) {
	cllctr := c.with()

	gauge, ok := cllctr.(prometheus.Gauge)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a %T)", cllctr, gauge))
	}
	gauge.Set(value)

	c.newValue(gauge)
}

func (c *Collector) Observe(value float64) {
	cllctr := c.with()

	observer, ok := cllctr.(prometheus.Histogram)
	if !ok {
		panic(fmt.Sprintf("invalid prometheus collector %T (not a %T)", cllctr, observer))
	}
	observer.Observe(value)

	c.newValue(observer)
}

func (c *Collector) newValue(value prometheus.Collector) {
	c.Values <- NewCollector(
		c.ID,
		"",
		c.Labels,
		value,
		c.Delete,
		nil,
	)
}

func NewCollectorID(name string, labels prometheus.Labels) string {
	var labelNamesValues []string
	for name, value := range labels {
		labelNamesValues = append(labelNamesValues, name, value)
	}
	sort.Strings(labelNamesValues)
	return name + ":" + strings.Join(labelNamesValues, "|")
}

func toLabels(labelValues ...string) prometheus.Labels {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}

	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}

	return labels
}
