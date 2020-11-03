package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Buckets []float64 `description:"Buckets for latency metrics." json:"buckets,omitempty" toml:"buckets,omitempty" yaml:"buckets,omitempty" export:"true"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Buckets: prometheus.DefBuckets,
	}
}
