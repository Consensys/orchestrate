package app

import (
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/log"
	metricsregister "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/registry"
)

type Config struct {
	HTTP    *HTTP
	GRPC    *GRPC
	Watcher *configwatcher.Config
	Log     *traefiktypes.TraefikLog
	Metrics *metricsregister.Config
}

type HTTP struct {
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

type GRPC struct {
	EntryPoint *traefikstatic.EntryPoint
	Static     *grpcstatic.Configuration
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		HTTP: &HTTP{
			EntryPoints: http.NewEPsConfig(vipr),
		},
		GRPC: &GRPC{
			EntryPoint: grpc.NewConfig(vipr),
		},
		Watcher: configwatcher.NewConfig(vipr),
		Log:     log.NewConfig(vipr),
		Metrics: metricsregister.NewConfig(vipr),
	}
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}
