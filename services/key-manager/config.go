package keymanager

import (
	"github.com/ConsenSys/orchestrate/pkg/hashicorp"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	httpmetrics "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/metrics"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	metricregistry "github.com/ConsenSys/orchestrate/pkg/toolkit/app/metrics/registry"
	tcpmetrics "github.com/ConsenSys/orchestrate/pkg/toolkit/tcp/metrics"
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
