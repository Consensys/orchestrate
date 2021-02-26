package static

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/router"
)

type Builder struct {
	routers map[string]*router.Router
}

func NewBuilder(routers map[string]*router.Router) *Builder {
	return &Builder{
		routers: routers,
	}
}

func (b *Builder) Build(ctx context.Context, entryPointNames []string, configuration interface{}) (map[string]*router.Router, error) {
	return b.routers, nil
}
