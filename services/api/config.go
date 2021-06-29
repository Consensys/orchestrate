package api

import (
	broker "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app"
	authjwt "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/jwt"
	authkey "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	httpmetrics "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/metrics"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	metricregistry "github.com/ConsenSys/orchestrate/pkg/toolkit/app/metrics/registry"
	tcpmetrics "github.com/ConsenSys/orchestrate/pkg/toolkit/tcp/metrics"
	"github.com/ConsenSys/orchestrate/services/api/metrics"
	"github.com/ConsenSys/orchestrate/services/api/proxy"
	store "github.com/ConsenSys/orchestrate/services/api/store/multi"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Flags register flags for API
func Flags(f *pflag.FlagSet) {
	log.Flags(f)
	multitenancy.Flags(f)
	authjwt.Flags(f)
	authkey.Flags(f)
	broker.KafkaProducerFlags(f)
	broker.KafkaTopicTxSender(f)
	qkm.Flags(f)
	store.Flags(f)
	http.Flags(f)
	http.MetricFlags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, tcpmetrics.ModuleName, metrics.ModuleName)
	proxy.Flags(f)
}

type Config struct {
	App          *app.Config
	Store        *store.Config
	Multitenancy bool
	Proxy        *proxy.Config
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:          app.NewConfig(vipr),
		Store:        store.NewConfig(vipr),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
		Proxy:        proxy.NewConfig(),
	}
}
