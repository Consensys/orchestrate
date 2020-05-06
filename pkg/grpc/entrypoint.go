package grpc

import (
	"context"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"google.golang.org/grpc"
)

const (
	DefaultGRPCEntryPoint = "grpc"
)

type EntryPoint struct {
	cfg *traefikstatic.EntryPoint

	tcp       *tcp.EntryPoint
	forwarder *tcp.Forwarder

	builder server.Builder
	server  *grpc.Server
}

func NewEntryPoint(name string, ep *traefikstatic.EntryPoint, builder server.Builder, reg metrics.TCP) *EntryPoint {
	forwarder := tcp.NewForwarder(nil)
	rt := &tcp.Router{}
	rt.TCPForwarder(forwarder)

	if name == "" {
		name = DefaultGRPCEntryPoint
	}

	return &EntryPoint{
		cfg:       ep,
		tcp:       tcp.NewEntryPoint(name, ep, rt, reg),
		forwarder: forwarder,
		builder:   builder,
	}
}

func (ep *EntryPoint) Addr() string {
	return ep.tcp.Addr()
}

func (ep *EntryPoint) BuildServer(ctx context.Context, configuration interface{}) error {
	var err error
	ep.server, err = ep.builder.Build(ctx, ep.tcp.Name(), configuration)
	return err
}

func (ep *EntryPoint) ListenAndServe(ctx context.Context) error {
	utils.InParallel(
		func() { _ = ep.server.Serve(ep.forwarder) },
		func() { _ = ep.tcp.ListenAndServe(ctx) },
	)

	return nil
}

func (ep *EntryPoint) Shutdown(ctx context.Context) error {
	utils.InParallel(
		func() { ep.server.GracefulStop() },
		func() { _ = tcp.Shutdown(ctx, ep.tcp) },
	)
	return nil
}

func (ep *EntryPoint) Close() error {
	utils.InParallel(
		func() { _ = tcp.Close(ep.forwarder) },
		func() { _ = tcp.Close(ep.tcp) },
	)
	return nil
}
