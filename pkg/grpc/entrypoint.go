package grpc

import (
	"context"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/hashicorp/go-multierror"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp/metrics"
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

func NewEntryPoint(name string, ep *traefikstatic.EntryPoint, builder server.Builder, reg metrics.TPCMetrics) *EntryPoint {
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
	go func() {
		// next error can be ignored
		// because net.Error are catched at the tcp.Entrypoint level
		_ = ep.server.Serve(ep.forwarder)
	}()

	return ep.tcp.ListenAndServe(ctx)
}

func (ep *EntryPoint) Shutdown(ctx context.Context) error {
	gr := &multierror.Group{}
	gr.Go(func() error { ep.server.GracefulStop(); return nil })
	gr.Go(func() error { return tcp.Shutdown(ctx, ep.tcp) })
	return gr.Wait().ErrorOrNil()
}

func (ep *EntryPoint) Close() error {
	gr := &multierror.Group{}
	gr.Go(func() error { return tcp.Close(ep.tcp) })
	gr.Go(func() error { return tcp.Close(ep.forwarder) })
	return gr.Wait().ErrorOrNil()
}
