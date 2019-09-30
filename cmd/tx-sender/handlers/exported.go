package handlers

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-sender/handlers/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-sender/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-sender/handlers/sender"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/handlers/opentracing"
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
		// Initialize nonce manager
		func() { nonce.Init(ctx) },
		// Initialize producer
		func() { producer.Init(ctx) },
	)
}
