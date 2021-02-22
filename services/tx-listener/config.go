package txlistener

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/key"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	metricregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/registry"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	provider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/providers/chain-registry"
	txsentry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sentry"
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
