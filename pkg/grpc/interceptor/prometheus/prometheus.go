package grpcprometheus

import (
	"context"
	"fmt"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"google.golang.org/grpc"
)

type Builder struct {
	metrics *grpc_prometheus.ServerMetrics
}

func NewBuilder(
	labels map[string]string,
	buckets []float64,
	histogramLabels map[string]string,
) *Builder {
	var opts []grpc_prometheus.CounterOption
	if len(labels) > 0 {
		opts = append(opts, grpc_prometheus.WithConstLabels(prom.Labels(labels)))
	}

	metrics := grpc_prometheus.NewServerMetrics(opts...)
	var histogramOpts []grpc_prometheus.HistogramOption
	if len(buckets) > 0 {
		histogramOpts = append(histogramOpts, grpc_prometheus.WithHistogramBuckets(buckets))
	}

	if len(histogramLabels) > 0 {
		histogramOpts = append(histogramOpts, grpc_prometheus.WithHistogramConstLabels(prom.Labels(histogramLabels)))
	}

	if len(histogramOpts) > 0 {
		metrics.EnableHandlingTimeHistogram(histogramOpts...)
	}

	return &Builder{
		metrics: metrics,
	}
}

func (b *Builder) ServerMetrics() *grpc_prometheus.ServerMetrics {
	return b.metrics
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Prometheus)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}

	return b.metrics.UnaryServerInterceptor(), b.metrics.StreamServerInterceptor(), b.metrics.InitializeMetrics, nil
}
