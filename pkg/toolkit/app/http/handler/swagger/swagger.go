package swagger

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	swaggerui "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/handler/swagger/genstatic"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

type Builder struct{}

// NewBuilder returns a http.Handler builder based on runtime.Configuration
func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.Swagger)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	New(cfg.SpecsFile).Append(router)

	return router, nil
}

type Swagger struct {
	serveSpecs http.Handler
	serveUI    http.Handler
}

func New(specsFile string) *Swagger {
	return &Swagger{
		serveUI: http.FileServer(&assetfs.AssetFS{Asset: swaggerui.Asset, AssetDir: swaggerui.AssetDir, Prefix: "public/swagger-ui"}),
		serveSpecs: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			http.ServeFile(rw, req, specsFile)
		}),
	}
}

// Append add dashboard routes on a router
func (h *Swagger) Append(router *mux.Router) {
	router.Methods(http.MethodGet).
		Path("/swagger/swagger.json").
		Handler(h.serveSpecs)

	router.Methods(http.MethodGet).
		PathPrefix("/swagger").
		Handler(http.StripPrefix("/swagger/", h.serveUI))
}
