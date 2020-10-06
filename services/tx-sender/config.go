package txsender

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
}

const (
	MetricsURLViperKey = "tx-sender.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "TX_SENDER_METRICS_URL"
)
