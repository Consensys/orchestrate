package transactionscheduler

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
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/multi"
)

// Flags register flags for tx scheduler
func Flags(f *pflag.FlagSet) {
	// Register Kafka flags
	broker.InitKafkaFlags(f)
	broker.KafkaTopicTxCrafter(f)
	broker.KafkaTopicTxSender(f)

	// Internal API clients
	chnregclient.Flags(f)
	registryclient.ContractRegistryURL(f)

	multi.Flags(f)
	http.Flags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, tcpmetrics.ModuleName, metrics.ModuleName)
}

type Config struct {
	App          *app.Config
	Store        *multi.Config
	Multitenancy bool
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:          app.NewConfig(vipr),
		Store:        multi.NewConfig(vipr),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}
}
