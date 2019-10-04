package grpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestInit(t *testing.T) {
	var s *grpc.Server
	SetGlobalServer(s)
	assert.Nil(t, GlobalServer(), "Global should be reset to nil")

	// Init gRPC server
	AddEnhancers(
		func(srv *grpc.Server) *grpc.Server {
			helloworld.RegisterGreeterServer(srv, &DummyGreeterServer{})
			return srv
		},
	)
	Init(context.Background())
	assert.NotNil(t, server, "Server should have been set")

	go ListenAndServe()
	assert.NotNil(t, server, "Server should have been set")
	StopServer(context.Background())
}
