package txsigner

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
}

const (
	MetricsURLViperKey = "tx-signer.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "TX_SIGNER_METRICS_URL"
)
