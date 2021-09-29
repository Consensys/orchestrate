package multi

import (
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	promcol "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/collector"
	"github.com/prometheus/client_golang/prometheus"
)

type DynamicMulti struct {
	Collectors []func(ch chan<- *prometheus.Desc)

	mux           *sync.Mutex
	currentCfg    *DynamicConfig
	currentValues map[string]*promcol.Collector
	values        chan *promcol.Collector
}

func NewDynamicMulti(cfg *DynamicConfig) *DynamicMulti {
	cllctr := &DynamicMulti{
		mux:           &sync.Mutex{},
		currentValues: make(map[string]*promcol.Collector),
		values:        make(chan *promcol.Collector, 100),
	}

	if cfg == nil {
		cllctr.currentCfg = NewDynamicConfig(nil)
	} else {
		cllctr.currentCfg = cfg
	}

	go cllctr.listenValues()

	return cllctr
}

// Describe implements prometheus.Collector and simply calls
// the registered describer functions.
func (c *DynamicMulti) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.Collectors {
		desc(ch)
	}
}

// Collect current collectors and then removes collectors that have
// removed of the dynamic configuration
func (c *DynamicMulti) Collect(ch chan<- prometheus.Metric) {
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
		defer promcol.CollectorPool.Put(c.currentValues[key])

		c.currentValues[key].Delete()
		delete(c.currentValues, key)
	}
}

func (c *DynamicMulti) Switch(cfg *dynamic.Configuration) error {
	c.mux.Lock()
	c.currentCfg = NewDynamicConfig(cfg)
	c.mux.Unlock()
	return nil
}

func (c *DynamicMulti) Values() chan *promcol.Collector {
	return c.values
}

func (c *DynamicMulti) listenValues() {
	for value := range c.values {
		c.listenValue(value)
	}
}

func (c *DynamicMulti) listenValue(value *promcol.Collector) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if oldValue, ok := c.currentValues[value.ID]; ok {
		defer promcol.CollectorPool.Put(oldValue)
	}

	c.currentValues[value.ID] = value
}

// isOutdated checks whether the passed collector has labels that mark
// it as belonging to an outdated dynamic configuration
func (c *DynamicMulti) isOutdated(value *promcol.Collector) bool {
	labels := value.Labels
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
