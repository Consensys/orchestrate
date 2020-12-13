package txsigner

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	noncechecker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/nonce/checker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/metrics"
	metricregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/registry"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
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

// Flags register flags for tx sentry
func Flags(f *pflag.FlagSet) {
	broker.InitKafkaFlags(f)
	keymanager.Flags(f)
	broker.KafkaTopicTxCrafter(f)
	broker.KafkaTopicTxRecover(f)
	broker.KafkaTopicTxSigner(f)
	chnregclient.Flags(f)
	noncechecker.Flags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, tcpmetrics.ModuleName)
}

type Config struct {
	App                *app.Config
	GroupName          string
	ListenerTopic      string
	RecoverTopic       string
	CrafterTopic       string
	ChainRegistryURL   string
	CheckerMaxRecovery uint64
	BckOff             backoff.BackOff
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:                app.NewConfig(vipr),
		GroupName:          "group-dispatcher",
		ListenerTopic:      vipr.GetString(broker.TxSignerViperKey),
		RecoverTopic:       vipr.GetString(broker.TxRecoverViperKey),
		CrafterTopic:       vipr.GetString(broker.TxCrafterViperKey),
		ChainRegistryURL:   vipr.GetString(chnregclient.URLViperKey),
		CheckerMaxRecovery: vipr.GetUint64(noncechecker.MaxRecoveryViperKey),
		BckOff:             retryMessageBackOff(),
	}
}

func retryMessageBackOff() backoff.BackOff {
	bckOff := backoff.NewExponentialBackOff()
	bckOff.MaxInterval = time.Second * 15
	bckOff.MaxElapsedTime = time.Minute * 5
	return bckOff
}
