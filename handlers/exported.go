package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers/sender"
)

type serviceName string

// Init inialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			ctxWithValue := context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctxWithValue)
		},
		// Initialize sender
		func() { sender.Init(ctx) },
	)
}
