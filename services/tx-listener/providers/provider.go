package provider

import (
	"context"

	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
)

//go:generate mockgen -source=provider.go -destination=mock/provider.go -package=mock

// Provider defines methods of a provider.
type Provider interface {
	// Run starts the provider to provide configuration to the tx-listener
	// Canceling ctx stops the provider
	// Once context is canceled Run should not send any message into the configuration input
	Run(ctx context.Context, configInput chan<- *dynamic.Message) error
}
