package app

import (
	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	metricsregister "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/registry"
	"github.com/spf13/viper"
	traefikstatic "github.com/traefik/traefik/v2/pkg/config/static"
	traefiktypes "github.com/traefik/traefik/v2/pkg/types"
)

type Config struct {
	HTTP    *HTTP
	Watcher *configwatcher.Config
	Log     *log.Config
	Metrics *metricsregister.Config
}

type HTTP struct {
	AccessLog    bool                             `description:"AccessLog enabled." json:"accessLog" toml:"accessLog" yaml:"accessLog" export:"true"`
	EntryPoints  traefikstatic.EntryPoints        `description:"Entry points definition." json:"entryPoints,omitempty" toml:"entryPoints,omitempty" yaml:"entryPoints,omitempty" export:"true"`
	HostResolver *traefiktypes.HostResolverConfig `description:"Enable CNAME Flattening." json:"hostResolver,omitempty" toml:"hostResolver,omitempty" yaml:"hostResolver,omitempty" label:"allowEmpty" export:"true"`
}

func (c *HTTP) TraefikStatic() *traefikstatic.Configuration {
	return &traefikstatic.Configuration{
		EntryPoints:  c.EntryPoints,
		HostResolver: c.HostResolver,
		API: &traefikstatic.API{
			Dashboard: true,
		},
	}
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		HTTP: &HTTP{
			EntryPoints: NewEPsConfig(vipr),
			AccessLog:   vipr.GetBool(accessLogEnabledKey),
		},
		Watcher: configwatcher.NewConfig(vipr),
		Log:     log.NewConfig(vipr),
		Metrics: metricsregister.NewConfig(vipr),
	}
}

func NewEPsConfig(vipr *viper.Viper) traefikstatic.EntryPoints {
	httpEp := &traefikstatic.EntryPoint{
		Address: url(vipr.GetString(hostnameViperKey), vipr.GetUint(httpPortViperKey)),
	}
	httpEp.SetDefaults()

	metricsEp := &traefikstatic.EntryPoint{
		Address: url(vipr.GetString(metricsHostnameViperKey), vipr.GetUint(metricsPortViperKey)),
	}
	metricsEp.SetDefaults()

	return traefikstatic.EntryPoints{
		http.DefaultHTTPAppEntryPoint: httpEp,
		http.DefaultMetricsEntryPoint: metricsEp,
	}
}
