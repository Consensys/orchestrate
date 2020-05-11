package client

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(txSchedulerURLViperKey, txSchedulerURLDefault)
	_ = viper.BindEnv(txSchedulerURLViperKey, txSchedulerURLEnv)
}

const (
	txSchedulerURLFlag     = "transaction-scheduler-url"
	txSchedulerURLViperKey = "transaction.scheduler.url"
	txSchedulerURLDefault  = "localhost:8081"
	txSchedulerURLEnv      = "TRANSACTION_SCHEDULER_URL"
)

// ChainRegistryURL register flag for the URL of the Chain Registry
func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL of the Transaction Scheduler HTTP endpoint. 
Environment variable: %q`, txSchedulerURLEnv)
	f.String(txSchedulerURLFlag, txSchedulerURLDefault, desc)
	_ = viper.BindPFlag(txSchedulerURLViperKey, f.Lookup(txSchedulerURLFlag))
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
		URL: vipr.GetString(txSchedulerURLViperKey),
	}
}
