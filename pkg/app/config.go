package app

import (
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	metrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
)

type Config struct {
	HTTP    *traefikstatic.Configuration
	GRPC    *GRPC
	Watcher *configwatcher.Config
	Metrics *metrics.Config
}

type GRPC struct {
	EntryPoint *traefikstatic.EntryPoint
	Static     *grpcstatic.Configuration
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		HTTP: http.NewConfig(vipr),
		GRPC: &GRPC{
			EntryPoint: grpc.NewConfig(vipr),
		},
		Watcher: configwatcher.NewConfig(vipr),
		Metrics: metrics.NewConfig(vipr),
	}
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}
