package transactionscheduler

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/multi"
)

const (
	defaultRetryIntervalFlag     = "default-retry-interval"
	defaultRetryIntervalViperKey = "retry-interval"
	defaultRetryIntervalDefault  = 45 * time.Second
	defaultRetryIntervalEnv      = "DEFAULT_RETRY_INTERVAL"
)

func init() {
	viper.SetDefault(defaultRetryIntervalViperKey, defaultRetryIntervalDefault)
	_ = viper.BindEnv(defaultRetryIntervalViperKey, defaultRetryIntervalEnv)
}

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

	defaultRetryIntervalDesc := fmt.Sprintf(`Amount of time a pending schedule needs to be considered for retry. Environment variable: %q`, defaultRetryIntervalEnv)
	f.Duration(defaultRetryIntervalFlag, defaultRetryIntervalDefault, defaultRetryIntervalDesc)
	_ = viper.BindPFlag(defaultRetryIntervalViperKey, f.Lookup(defaultRetryIntervalFlag))
}

type Config struct {
	App           *app.Config
	Store         *multi.Config
	RetryInterval time.Duration
	Multitenancy  bool
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:           app.NewConfig(vipr),
		Store:         multi.NewConfig(vipr),
		RetryInterval: vipr.GetDuration(defaultRetryIntervalViperKey),
		Multitenancy:  viper.GetBool(multitenancy.EnabledViperKey),
	}
}
