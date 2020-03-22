package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

//go:generate mockgen -source=interceptor.go -destination=mock/mock.go -package=mock

type Builder interface {
	Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error)
}
