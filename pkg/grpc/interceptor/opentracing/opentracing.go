package grpcopentracing

import (
	"context"
	"fmt"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"google.golang.org/grpc"
)

type Builder struct {
	tracer opentracing.Tracer
}

func NewBuilder(tracer opentracing.Tracer) *Builder {
	return &Builder{
		tracer: tracer,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Tracing)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}

	opts := []grpc_opentracing.Option{grpc_opentracing.WithTracer(b.tracer)}

	// TODO: Add Option for TraceHeaderName when grpc-ecosystem/go-grpc-middleware release next version>1.2.0
	// if cfg.TraceHeaderName != "" {
	// 	opts = append(opts, grpc_opentracing.WithTraceHeaderName(cfg.TraceHeaderName))
	// }

	return grpc_opentracing.UnaryServerInterceptor(opts...), grpc_opentracing.StreamServerInterceptor(opts...), nil, nil
}
