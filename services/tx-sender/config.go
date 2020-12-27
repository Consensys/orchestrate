package txsender

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/redis"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/metrics"
	metricregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/registry"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tcp/metrics"
	chnregclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

func init() {
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)

	viper.SetDefault(NonceMaxRecoveryViperKey, nonceMaxRecoveryDefault)
	_ = viper.BindEnv(NonceMaxRecoveryViperKey, nonceMaxRecoveryEnv)

	viper.SetDefault(nonceManagerTypeViperKey, nonceManagerTypeDefault)
	_ = viper.BindEnv(nonceManagerTypeViperKey, nonceManagerTypeEnv)

	viper.SetDefault(NonceManagerExpirationViperKey, nonceManagerExpirationDefault)
	_ = viper.BindEnv(NonceManagerExpirationViperKey, nonceManagerExpirationEnv)
}

const (
	MetricsURLViperKey = "tx-sender.metrics.url"
	metricsURLDefault  = "localhost:8082"
	metricsURLEnv      = "TX_SENDER_METRICS_URL"
)

const (
	nonceMaxRecoveryFlag     = "nonce-max-recovery"
	NonceMaxRecoveryViperKey = "nonce.max.recovery"
	nonceMaxRecoveryDefault  = 5
	nonceMaxRecoveryEnv      = "NONCE_MAX_RECOVERY"
)

const (
	nonceManagerTypeFlag     = "nonce-manager-type"
	nonceManagerTypeViperKey = "nonce.manager.type"
	nonceManagerTypeDefault  = "redis"
	nonceManagerTypeEnv      = "NONCE_MANAGER_TYPE"

	NonceManagerTypeInMemory = "in-memory"
	NonceManagerTypeRedis    = "redis"
)

const (
	nonceManagerExpirationFlag     = "nonce-manager-expiration"
	NonceManagerExpirationViperKey = "nonce.manager.expiration"
	nonceManagerExpirationDefault  = 5 * time.Minute
	nonceManagerExpirationEnv      = "NONCE_MANAGER_EXPIRATION"
)

// Flags register flags for tx sentry
func Flags(f *pflag.FlagSet) {
	broker.InitKafkaFlags(f)
	keymanager.Flags(f)
	broker.KafkaTopicTxSender(f)
	broker.KafkaTopicTxRecover(f)
	chnregclient.Flags(f)
	MaxRecovery(f)
	NonceManagerType(f)
	NonceManagerExpirationFlag(f)
	redis.Flags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, tcpmetrics.ModuleName)
	httputils.MetricFlags(f)
}

// MaxRecovery register a flag for Redis server MaxRecovery
func MaxRecovery(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of times tx-sender tries to recover an envelope with invalid nonce before outputing an error
Environment variable: %q`, nonceMaxRecoveryEnv)
	f.Int(nonceMaxRecoveryFlag, nonceMaxRecoveryDefault, desc)
	_ = viper.BindPFlag(NonceMaxRecoveryViperKey, f.Lookup(nonceMaxRecoveryFlag))
}

// Type register flag for Nonce Cooldown
func NonceManagerType(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Nonce (one of %q)
Environment variable: %q`, []string{NonceManagerTypeInMemory, NonceManagerTypeRedis}, nonceManagerTypeEnv)
	f.String(nonceManagerTypeFlag, nonceManagerTypeDefault, desc)
	_ = viper.BindPFlag(nonceManagerTypeViperKey, f.Lookup(nonceManagerTypeFlag))
}

// ExpirationFlag register a flag for Redis expiration
func NonceManagerExpirationFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`NonceManager values expiration time.
Environment variable: %q`, nonceManagerExpirationEnv)
	f.Duration(nonceManagerExpirationFlag, nonceManagerExpirationDefault, desc)
	_ = viper.BindPFlag(NonceManagerExpirationViperKey, f.Lookup(nonceManagerExpirationFlag))
}

type Config struct {
	App                    *app.Config
	GroupName              string
	RecoverTopic           string
	SenderTopic            string
	ChainRegistryURL       string
	BckOff                 backoff.BackOff
	NonceMaxRecovery       uint64
	NonceManagerType       string
	RedisCfg               *redis.Config
	NonceManagerExpiration time.Duration
}

func NewConfig(vipr *viper.Viper) *Config {
	redisCfg := redis.NewConfig(vipr)
	redisCfg.Expiration = int(vipr.GetDuration(NonceManagerExpirationViperKey).Milliseconds())

	return &Config{
		App:                    app.NewConfig(vipr),
		GroupName:              "group-dispatcher",
		RecoverTopic:           vipr.GetString(broker.TxRecoverViperKey),
		SenderTopic:            vipr.GetString(broker.TxSenderViperKey),
		ChainRegistryURL:       vipr.GetString(chnregclient.URLViperKey),
		NonceMaxRecovery:       vipr.GetUint64(NonceMaxRecoveryViperKey),
		BckOff:                 retryMessageBackOff(),
		NonceManagerType:       viper.GetString(nonceManagerTypeViperKey),
		NonceManagerExpiration: vipr.GetDuration(NonceManagerExpirationViperKey),
		RedisCfg:               redisCfg,
	}
}

func retryMessageBackOff() backoff.BackOff {
	bckOff := backoff.NewExponentialBackOff()
	bckOff.MaxInterval = time.Second * 15
	bckOff.MaxElapsedTime = time.Minute * 5
	return bckOff
}
