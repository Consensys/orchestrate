package api

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/metrics"
	metricregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/multi"
	keymanagerclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

// Flags register flags for API
func Flags(f *pflag.FlagSet) {
	// Register Kafka flags
	broker.InitKafkaFlags(f)
	broker.KafkaTopicTxSender(f)

	// Internal API clients
	keymanagerclient.Flags(f)

	multi.Flags(f)
	http.Flags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, tcpmetrics.ModuleName, metrics.ModuleName)

	proxy.Flags(f)
}

type Config struct {
	App          *app.Config
	Store        *multi.Config
	Multitenancy bool
	Proxy        *proxy.Config
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:          app.NewConfig(vipr),
		Store:        multi.NewConfig(vipr),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
		Proxy:        proxy.NewConfig(),
	}
}
