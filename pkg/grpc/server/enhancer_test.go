package grpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type DummyGreeterServer struct{}

// SayHello implements helloworld.GreeterServer
func (s *DummyGreeterServer) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + req.Name}, nil
}

func TestApplyEnhancers(t *testing.T) {
	srv := grpc.NewServer()
	ApplyEnhancers(
		srv,
		func(srv *grpc.Server) *grpc.Server {
			helloworld.RegisterGreeterServer(srv, &DummyGreeterServer{})
			return srv
		},
	)

	infos := srv.GetServiceInfo()
	assert.Len(t, infos, 1, "Service should have been registered")
}
