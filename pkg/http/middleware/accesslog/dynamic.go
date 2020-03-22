package accesslog

import (
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

func AddDynamicConfig(cfg *dynamic.Configuration, midName string, staticCfg *traefikstatic.Configuration) {
	// Access Log Middleware
	logFormat := ""
	if staticCfg.Log != nil {
		logFormat = staticCfg.Log.Format
	}

	cfg.HTTP.Middlewares[midName] = &dynamic.Middleware{
		AccessLog: &dynamic.AccessLog{
			Format: logFormat,
		},
	}

}
