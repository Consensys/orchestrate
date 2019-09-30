package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-listener/handlers/enricher"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-listener/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-listener/handlers/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/handlers/opentracing"
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
