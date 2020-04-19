package multi

import (
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/prometheus"
)

type Config struct {
	Prometheus *prometheus.Config
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Prometheus: prometheus.NewConfig(vipr),
	}
}
