package grpcmetrics

import (
	"context"
	"fmt"
	"time"

	kitmetrics "github.com/go-kit/kit/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/grpcutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Builder struct {
	registry metrics.GRPCServer
}

func NewBuilder(registry metrics.GRPCServer) *Builder {
	return &Builder{
		registry: registry,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Metrics)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}
	return UnaryServerInterceptor(b.registry), StreamServerInterceptor(b.registry), nil, nil
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryServerInterceptor(registry metrics.GRPCServer) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service, method := grpcutil.SplitMethodName(info.FullMethod)
		baseLabels := []string{
			"tenant_id", multitenancy.TenantIDFromContext(ctx),
			"service", service,
			"type", grpcutil.Unary,
			"method", method,
		}

		registry.StartedCounter().With(baseLabels...).Add(1)
		start := time.Now()

		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)

		handledLabels := append(
			baseLabels,
			"code", st.Code().String(),
		)
		registry.HandledCounter().With(handledLabels...).Add(1)
		registry.HandledDurationHistogram().With(handledLabels...).Observe(time.Since(start).Seconds())

		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(registry metrics.GRPCServer) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service, method := grpcutil.SplitMethodName(info.FullMethod)
		baseLabels := []string{
			"tenant_id", multitenancy.TenantIDFromContext(ss.Context()),
			"service", service,
			"type", grpcutil.TypeFromStreamServerInfo(info),
			"method", method,
		}

		registry.StartedCounter().With(baseLabels...).Add(1)
		start := time.Now()

		err := handler(
			srv,
			&metricsServerStream{
				ServerStream:       ss,
				receivedMsgCounter: registry.StreamMsgReceivedCounter().With(baseLabels...),
				sentMsgCounter:     registry.StreamMsgSentCounter().With(baseLabels...),
			},
		)

		st, _ := status.FromError(err)
		handledLabels := append(
			baseLabels,
			"code", st.Code().String(),
		)
		registry.HandledCounter().With(handledLabels...).Add(1)
		registry.HandledDurationHistogram().With(handledLabels...).Observe(time.Since(start).Seconds())

		return err
	}
}

type metricsServerStream struct {
	grpc.ServerStream
	receivedMsgCounter kitmetrics.Counter
	sentMsgCounter     kitmetrics.Counter
}

func (s *metricsServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.sentMsgCounter.Add(1)
	}

	return err
}

func (s *metricsServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.receivedMsgCounter.Add(1)
	}

	return err
}
