package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-nonce/handlers/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-nonce/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/handlers/opentracing"
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
		// Initialize Nonce manager
		func() {
			nonce.Init(ctx)
		},
		// Initialize PrepareMsg
		func() {
			producer.Init(ctx)
		},
	)
}
