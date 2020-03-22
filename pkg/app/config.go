package app

import (
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

type Config struct {
	HTTP    *traefikstatic.Configuration
	GRPC    *traefikstatic.EntryPoint
	Watcher *configwatcher.Config
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		HTTP:    http.NewConfig(vipr),
		GRPC:    grpc.NewConfig(vipr),
		Watcher: configwatcher.NewConfig(vipr),
	}
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}
