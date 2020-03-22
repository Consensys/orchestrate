package log

import (
	"strings"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
)

func ConfigureLogger(staticCfg *traefikstatic.Configuration) {
	if staticCfg.Log != nil && staticCfg.Log.Level != "" {
		level, err := logrus.ParseLevel(strings.ToLower(staticCfg.Log.Level))
		if err != nil {
			log.WithoutContext().WithError(err).Errorf("Error getting level: %v", err)
		}
		log.SetLevel(level)
	}
}
