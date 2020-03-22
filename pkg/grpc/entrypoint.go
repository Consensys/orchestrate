package grpc

import (
	"context"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
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

	server *grpc.Server
}

func NewEntryPoint(ep *traefikstatic.EntryPoint, server *grpc.Server) *EntryPoint {
	forwarder := tcp.NewForwarder(nil)
	rt := &tcp.Router{}
	rt.TCPForwarder(forwarder)

	return &EntryPoint{
		cfg:       ep,
		tcp:       tcp.NewEntryPoint(ep, rt),
		forwarder: forwarder,
		server:    server,
	}
}

func (ep *EntryPoint) Addr() string {
	return ep.tcp.Addr()
}

func (ep *EntryPoint) with(ctx context.Context) context.Context {
	return log.With(ctx, log.Str("entrypoint", DefaultGRPCEntryPoint))
}

func (ep *EntryPoint) ListenAndServe(ctx context.Context) error {
	l, err := tcp.Listen(ep.cfg.Address)
	if err != nil {
		return err
	}

	utils.InParallel(
		func() { _ = ep.server.Serve(ep.forwarder) },
		func() { _ = ep.tcp.Serve(ep.with(ctx), l) },
	)

	return nil
}

func (ep *EntryPoint) Shutdown(ctx context.Context) error {
	utils.InParallel(
		func() { ep.server.GracefulStop() },
		func() { _ = tcp.Shutdown(ep.with(ctx), ep.tcp) },
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
