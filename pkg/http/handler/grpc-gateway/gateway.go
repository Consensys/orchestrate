package grpcgateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type Builder struct {
	opts         []runtime.ServeMuxOption
	registrators []func(ctx context.Context, mux *runtime.ServeMux) error
}

func NewBuilder(opts []runtime.ServeMuxOption, registrators []func(ctx context.Context, mux *runtime.ServeMux) error) *Builder {
	return &Builder{
		opts:         opts,
		registrators: registrators,
	}
}

func (b *Builder) Build(ctx context.Context, _ string, _ interface{}, _ func(*http.Response) error) (http.Handler, error) {
	mux := runtime.NewServeMux(b.opts...)
	for _, registrator := range b.registrators {
		err := registrator(ctx, mux)
		if err != nil {
			return nil, err
		}
	}
	return mux, nil
}
