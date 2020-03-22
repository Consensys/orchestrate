package provider

import (
	"context"
)

// Message holds configuration information exchanged between
// a provider and
type Message interface {
	ProviderName() string
	Configuration() interface{}
}

// Provider defines methods of a provider.
type Provider interface {
	// Provide allows the provider to provide messages to a channel
	Provide(ctx context.Context, msgs chan<- Message) error
}
