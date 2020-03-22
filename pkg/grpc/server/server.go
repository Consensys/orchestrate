package server

import (
	"context"

	"google.golang.org/grpc"
)

//go:generate mockgen -source=server.go -destination=mock/mock.go -package=mock

type Builder interface {
	Build(ctx context.Context, name string, configuration interface{}) (*grpc.Server, error)
}

type OptionsBuilder interface {
	Build(ctx context.Context, name string, configuration interface{}) ([]grpc.ServerOption, error)
}
