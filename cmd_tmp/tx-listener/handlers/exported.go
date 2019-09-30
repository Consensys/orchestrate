package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/enricher"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/store"
)

type serviceName string

// Init initialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger
		func() {
			ctx = context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctx)
		},
		// Initialize store
		func() {
			store.Init(ctx)
		},
		// Initialize enricher
		func() {
			enricher.Init(ctx)
		},
		// Initialize producer
		func() {
			producer.Init(ctx)
		},
	)
}
