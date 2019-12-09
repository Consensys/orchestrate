package chainregistry

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/roundrobin"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/genstatic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/server"
)

var (
	startOnce = &sync.Once{}
)

// Run application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		staticConfig := chainregistry.NewConfig()

		// Configure logger
		configureLogging(staticConfig)

		staticConfig.API.DashboardAssets = &assetfs.AssetFS{Asset: genstatic.Asset, AssetDir: genstatic.AssetDir, Prefix: "static"}

		// Prepare configuration
		http.DefaultTransport.(*http.Transport).Proxy = http.ProxyFromEnvironment

		if err := roundrobin.SetDefaultWeight(0); err != nil {
			log.WithoutContext().WithError(err).Fatal("could not set round robin default weight")
		}

		staticConfig.SetEffectiveConfiguration()
		if err := staticConfig.ValidateConfiguration(); err != nil {
			log.WithoutContext().WithError(err).Fatal("invalid configuration")
		}

		jsonConf, err := json.Marshal(staticConfig)
		if err != nil {
			log.WithoutContext().WithError(err).Fatalf("could not marshal static configuration: %#v", staticConfig)
		} else {
			log.WithoutContext().Infof("static configuration loaded %s", string(jsonConf))
		}

		// Initialize server
		server.SetGlobalStaticConfig(staticConfig)
		server.Init(ctx)

		// Start server
		server.GlobalServer().Start(ctx)

		// Wait for server to properly close
		server.GlobalServer().Wait()

		log.WithoutContext().Info("Shutting down")
	})
}

func configureLogging(staticConfig *static.Configuration) {
	if staticConfig.Log != nil && staticConfig.Log.Level != "" {
		level, err := logrus.ParseLevel(strings.ToLower(staticConfig.Log.Level))
		if err != nil {
			log.WithoutContext().WithError(err).Errorf("Error getting level: %v", err)
		}
		log.SetLevel(level)
	}
}
