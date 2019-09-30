package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/vault"
)

type serviceName string

// Init inialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			ctx = context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctx)
		},
		// Initialize Vault
		func() { vault.Init(ctx) },
		// Initialize Producer
		func() { producer.Init(ctx) },
	)
}
