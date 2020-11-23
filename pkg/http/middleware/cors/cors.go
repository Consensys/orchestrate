package cors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
)

type Builder struct{}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.Cors)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	m := cors.New(cors.Options{
		AllowedOrigins:     cfg.AllowedOrigins,
		AllowedMethods:     cfg.AllowedMethods,
		AllowedHeaders:     cfg.AllowedHeaders,
		ExposedHeaders:     cfg.ExposedHeaders,
		MaxAge:             cfg.MaxAge,
		AllowCredentials:   cfg.AllowCredentials,
		OptionsPassthrough: cfg.OptionsPassthrough,
	})

	return m.Handler, nil, nil
}
