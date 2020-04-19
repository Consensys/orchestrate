package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

type Config struct {
	TCP  *TCPConfig        `description:"Entry points metrics configuration" json:"entryPoint,omitempty" toml:"entryPoint,omitempty" yaml:"entryPoint,omitempty" export:"true"`
	HTTP *HTTPConfig       `description:"HTTP metrics configuration" json:"http,omitempty" toml:"http,omitempty" yaml:"http,omitempty" export:"true"`
	GRPC *GRPCServerConfig `description:"GRPC metrics configuration" json:"grpc,omitempty" toml:"grpc,omitempty" yaml:"grpc,omitempty" export:"true"`
}

type TCPConfig struct {
	Buckets []float64 `description:"Buckets for latency metrics." json:"buckets,omitempty" toml:"buckets,omitempty" yaml:"buckets,omitempty" export:"true"`
}

type HTTPConfig struct {
	Buckets []float64 `description:"Buckets for latency metrics." json:"buckets,omitempty" toml:"buckets,omitempty" yaml:"buckets,omitempty" export:"true"`
}

type GRPCServerConfig struct {
	Buckets []float64 `description:"Buckets for latency metrics." json:"buckets,omitempty" toml:"buckets,omitempty" yaml:"buckets,omitempty" export:"true"`
}

func NewConfig(_ *viper.Viper) *Config {
	return &Config{
		TCP: &TCPConfig{
			Buckets: prometheus.DefBuckets,
		},
		HTTP: &HTTPConfig{
			Buckets: prometheus.DefBuckets,
		},
		GRPC: &GRPCServerConfig{
			Buckets: prometheus.DefBuckets,
		},
	}
}
