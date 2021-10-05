package txlistener

import (
	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app"
	authkey "github.com/consensys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	metricregistry "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/registry"
	tcpmetrics "github.com/consensys/orchestrate/pkg/toolkit/tcp/metrics"
	provider "github.com/consensys/orchestrate/services/tx-listener/providers/chain-registry"
	txsentry "github.com/consensys/orchestrate/services/tx-sentry"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	MetricsURLViperKey = "tx-listener.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "TX_LISTENER_METRICS_URL"
)

func init() {
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
}

// Flags register flags for API
func Flags(f *pflag.FlagSet) {
	log.Flags(f)
	authkey.Flags(f)
	broker.KafkaProducerFlags(f)
	broker.KafkaTopicTxDecoded(f)
	app.MetricFlags(f)
	metricregistry.Flags(f, tcpmetrics.ModuleName)
	txsentry.Flags(f)
	provider.Flags(f)
	orchestrateclient.Flags(f)
}
