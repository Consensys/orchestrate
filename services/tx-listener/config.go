package txlistener

import (
	broker "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	orchestrateclient "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	authkey "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	metricregistry "github.com/ConsenSys/orchestrate/pkg/toolkit/app/metrics/registry"
	tcpmetrics "github.com/ConsenSys/orchestrate/pkg/toolkit/tcp/metrics"
	provider "github.com/ConsenSys/orchestrate/services/tx-listener/providers/chain-registry"
	txsentry "github.com/ConsenSys/orchestrate/services/tx-sentry"
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
	http.MetricFlags(f)
	metricregistry.Flags(f, tcpmetrics.ModuleName)
	txsentry.Flags(f)
	provider.Flags(f)
	orchestrateclient.Flags(f)
}
