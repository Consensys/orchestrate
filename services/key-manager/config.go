package keymanager

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/metrics"
	metricregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/registry"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp/metrics"
	store "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store/multi"
)

// Flags register flags for tx scheduler
func Flags(f *pflag.FlagSet) {
	store.Flags(f)
	http.Flags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, tcpmetrics.ModuleName)
}

type Config struct {
	App   *app.Config
	Store *store.Config
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:   app.NewConfig(vipr),
		Store: store.NewConfig(vipr),
	}
}
