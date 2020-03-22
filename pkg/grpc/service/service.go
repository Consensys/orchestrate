package service

import (
	"context"

	"google.golang.org/grpc"
)

//go:generate mockgen -source=service.go  -destination=mock/mock.go -package=mock

type Builder interface {
	Build(ctx context.Context, name string, configuration interface{}) (func(srv *grpc.Server), error)
}
