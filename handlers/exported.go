package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/producer"
)

type serviceName string

// Init initialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			ctx = context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctx)
		},
		// Initialize decoder
		func() {
			decoder.Init(ctx)
		},
		// Initialize Producer
		func() {
			producer.Init(ctx)
		},
	)
}
