package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(TxSchedulerURLViperKey, TxSchedulerURLDefault)
	_ = viper.BindEnv(TxSchedulerURLViperKey, txSchedulerURLEnv)
	viper.SetDefault(TxSchedulerMetricsURLViperKey, TxSchedulerMetricsURLDefault)
	_ = viper.BindEnv(TxSchedulerMetricsURLViperKey, txSchedulerMetricsURLEnv)
}

const (
	txSchedulerURLFlag     = "transaction-scheduler-url"
	TxSchedulerURLViperKey = "transaction.scheduler.url"
	TxSchedulerURLDefault  = "localhost:8081"
	txSchedulerURLEnv      = "TRANSACTION_SCHEDULER_URL"
)

const (
	txSchedulerMetricsURLFlag     = "transaction.scheduler-metrics-url"
	TxSchedulerMetricsURLViperKey = "transaction.scheduler.metrics.url"
	TxSchedulerMetricsURLDefault  = "localhost:8082"
	txSchedulerMetricsURLEnv      = "TRANSACTION_SCHEDULER_METRICS_URL"
)

// ChainRegistryURL register flag for the URL of the Chain Registry
func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Transaction Scheduler HTTP endpoint. 
Environment variable: %q`, txSchedulerURLEnv)
	f.String(txSchedulerURLFlag, TxSchedulerURLDefault, desc)
	_ = viper.BindPFlag(TxSchedulerURLViperKey, f.Lookup(txSchedulerURLFlag))
}

// ChainRegistryURL register flag for the URL of the Chain Registry
func MetricsURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Transaction Scheduler Metrics endpoint.
Environment variable: %q`, txSchedulerMetricsURLEnv)
	f.String(txSchedulerMetricsURLFlag, TxSchedulerMetricsURLDefault, desc)
	_ = viper.BindPFlag(TxSchedulerMetricsURLViperKey, f.Lookup(txSchedulerMetricsURLFlag))
}
func Flags(f *pflag.FlagSet) {
	URL(f)
}

type Config struct {
	URL string
}

func NewConfig(url string) *Config {
	return &Config{
		URL: url,
	}
}

func NewConfigFromViper(vipr *viper.Viper) *Config {
	return &Config{
		URL: vipr.GetString(TxSchedulerURLViperKey),
	}
}
