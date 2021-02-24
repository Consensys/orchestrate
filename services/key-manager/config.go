package keymanager

import (
	"github.com/ConsenSys/orchestrate/pkg/app"
	"github.com/ConsenSys/orchestrate/pkg/hashicorp"
	"github.com/ConsenSys/orchestrate/pkg/http"
	httpmetrics "github.com/ConsenSys/orchestrate/pkg/http/metrics"
	"github.com/ConsenSys/orchestrate/pkg/log"
	metricregistry "github.com/ConsenSys/orchestrate/pkg/metrics/registry"
	tcpmetrics "github.com/ConsenSys/orchestrate/pkg/tcp/metrics"
	"github.com/ConsenSys/orchestrate/services/key-manager/store"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Flags register flags
func Flags(f *pflag.FlagSet) {
	log.Flags(f)
	store.Flags(f)
	http.Flags(f)
	http.MetricFlags(f)
	hashicorp.Flags(f)
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
