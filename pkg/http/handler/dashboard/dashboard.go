package dashboard

import (
	"context"
	"fmt"
	"net/http"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/runtime"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dashboard/genstatic"
)

type Builder struct {
	staticCfg *traefikstatic.Configuration
}

// NewBuilder returns a http.Handler builder based on runtime.Configuration
func NewBuilder(cfg *traefikstatic.Configuration) *Builder {
	return &Builder{
		staticCfg: cfg,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	if b.staticCfg == nil || b.staticCfg.API == nil {
		return nil, fmt.Errorf("dashboard is not enabled (consider updating static configuration)")
	}

	cfg, ok := configuration.(*runtime.Infos)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	return New(b.staticCfg, cfg), nil
}

// New returns a Handler defined by staticConfig, and if provided, by runtimeConfig.
// It finishes populating the information provided in the runtimeConfig.
func New(staticCfg *traefikstatic.Configuration, infos *runtime.Infos) http.Handler {
	if infos == nil {
		infos = &runtime.Infos{}
	}

	router := mux.NewRouter()
	NewOverview(staticCfg, infos).Append(router)
	NewEntryPoint(staticCfg).Append(router)
	NewHTTP(infos).Append(router)

	if staticCfg.API != nil && staticCfg.API.DashboardAssets != nil {
		NewUI(http.FileServer(staticCfg.API.DashboardAssets)).Append(router)
	} else {
		NewUI(http.FileServer(&assetfs.AssetFS{Asset: genstatic.Asset, AssetDir: genstatic.AssetDir, Prefix: "static"})).Append(router)
	}

	if staticCfg.API != nil && staticCfg.API.Debug {
		NewDebug().Append(router)
	}

	return router
}
