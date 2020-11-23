package grpclogrus

import (
	"context"
	"fmt"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/config/static"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

type Builder struct {
	logger *logrus.Logger
}

func NewBuilder(logger *logrus.Logger, fields logrus.Fields) *Builder {
	grpclog.SetLoggerV2(
		&Entry{
			Entry: logrus.WithFields(fields),
		},
	)

	return &Builder{
		logger: logger,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Logrus)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}

	fields := logrus.Fields{}
	for k, v := range cfg.Fields {
		fields[k] = v
	}

	return UnaryServerInterceptor(logrus.NewEntry(b.logger).WithFields(fields)), StreamServerInterceptor(logrus.NewEntry(b.logger)), nil, nil
}

// UnaryServerInterceptor returns a grpc unary server interceptor (middleware) that allows
// to intercept internal errors
func UnaryServerInterceptor(entry *logrus.Entry) grpc.UnaryServerInterceptor {
	return grpc_logrus.UnaryServerInterceptor(entry, grpc_logrus.WithLevels(CodeToLevel))
}

// StreamServerInterceptor returns a grpc streaming server interceptor for panic recovery.
func StreamServerInterceptor(entry *logrus.Entry) grpc.StreamServerInterceptor {
	return grpc_logrus.StreamServerInterceptor(entry, grpc_logrus.WithLevels(CodeToLevel))
}

func CodeToLevel(code codes.Code) logrus.Level {
	switch code {
	case codes.OK:
		return logrus.DebugLevel
	case codes.NotFound:
		return logrus.DebugLevel
	default:
		return grpc_logrus.DefaultCodeToLevel(code)
	}
}
