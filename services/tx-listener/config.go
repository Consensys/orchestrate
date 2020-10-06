package txlistener

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
}

const (
	MetricsURLViperKey = "tx-listener.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "TX_LISTENER_METRICS_URL"
)
