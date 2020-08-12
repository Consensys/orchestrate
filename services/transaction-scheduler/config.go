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
	sentryRefreshIntervalFlag     = "tx-sentry-refresh-interval"
	sentryRefreshIntervalViperKey = "tx-sentry.refresh-interval"
	sentryRefreshIntervalDefault  = 10 * time.Second
	sentryRefreshIntervalEnv      = "TX_SENTRY_REFRESH_INTERVAL"

	sentryPendingDurationFlag     = "tx-sentry-pending-duration"
	sentryPendingDurationViperKey = "tx-sentry.pending-duration"
	sentryPendingDurationDefault  = 2 * time.Minute
	sentryPendingDurationEnv      = "TX_SENTRY_PENDING_DURATION"
)

func init() {
	viper.SetDefault(sentryRefreshIntervalViperKey, sentryRefreshIntervalDefault)
	_ = viper.BindEnv(sentryRefreshIntervalViperKey, sentryRefreshIntervalEnv)

	viper.SetDefault(sentryPendingDurationViperKey, sentryPendingDurationDefault)
	_ = viper.BindEnv(sentryPendingDurationViperKey, sentryPendingDurationEnv)
}

// TxSchedulerFlags register flags for tx scheduler
func TxSchedulerFlags(f *pflag.FlagSet) {
	// Register Kafka flags
	broker.InitKafkaFlags(f)
	broker.KafkaTopicTxCrafter(f)
	broker.KafkaTopicTxSender(f)

	// Internal API clients
	chnregclient.Flags(f)
	registryclient.ContractRegistryURL(f)

	multi.Flags(f)
	http.Flags(f)
}

// TxSentryFlags register flags for tx sentry
func TxSentryFlags(f *pflag.FlagSet) {
	refreshIntervalDesc := fmt.Sprintf(`Time interval for refreshing the list of schedules. Environment variable: %q`, sentryRefreshIntervalEnv)
	f.Duration(sentryRefreshIntervalFlag, sentryRefreshIntervalDefault, refreshIntervalDesc)
	_ = viper.BindPFlag(sentryRefreshIntervalViperKey, f.Lookup(sentryRefreshIntervalFlag))

	pendingDurationDesc := fmt.Sprintf(`Amount of time a pending schedule needs to be considered for retry. Environment variable: %q`, sentryPendingDurationEnv)
	f.Duration(sentryPendingDurationFlag, sentryPendingDurationDefault, pendingDurationDesc)
	_ = viper.BindPFlag(sentryPendingDurationViperKey, f.Lookup(sentryPendingDurationFlag))
}

type Config struct {
	App          *app.Config
	Store        *multi.Config
	Sentry       *SentryConfig
	Multitenancy bool
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:          app.NewConfig(vipr),
		Store:        multi.NewConfig(vipr),
		Sentry:       NewSentryConfig(vipr),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}
}

type SentryConfig struct {
	RefreshInterval time.Duration
	PendingDuration time.Duration
}

func NewSentryConfig(vipr *viper.Viper) *SentryConfig {
	return &SentryConfig{
		RefreshInterval: vipr.GetDuration(sentryRefreshIntervalViperKey),
		PendingDuration: vipr.GetDuration(sentryPendingDurationViperKey),
	}
}
