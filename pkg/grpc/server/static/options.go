package static

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"google.golang.org/grpc"
)

type OptionsBuilder struct{}

func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{}
}

func (b *OptionsBuilder) Build(ctx context.Context, name string, configuration interface{}) ([]grpc.ServerOption, error) {
	cfg, ok := configuration.(*static.Options)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	if cfg == nil {
		return nil, nil
	}

	var opts []grpc.ServerOption
	if cfg.ConnectionTimeout != 0 {
		opts = append(opts, grpc.ConnectionTimeout(cfg.ConnectionTimeout))
	}

	if cfg.HeaderTableSize != 0 {
		opts = append(opts, grpc.HeaderTableSize(cfg.HeaderTableSize))
	}

	if cfg.MaxConcurrentStreams != 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(cfg.MaxConcurrentStreams))
	}

	if cfg.MaxHeaderListSize != 0 {
		opts = append(opts, grpc.MaxHeaderListSize(cfg.MaxHeaderListSize))
	}

	if cfg.MaxRecvMsgSize != 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize))
	}

	if cfg.MaxSendMsgSize != 0 {
		opts = append(opts, grpc.MaxSendMsgSize(cfg.MaxSendMsgSize))
	}

	return opts, nil
}
