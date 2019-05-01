package store

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// EnvelopeStore is an interface for context store
type EnvelopeStore interface {
	// Store context envelope
	// For a envelope with same TxHash that was already store, it updates the envelope and return status
	Store(ctx context.Context, envelope *envelope.Envelope) (status string, at time.Time, err error)

	// Load context envelope by txHash
	LoadByTxHash(ctx context.Context, chainID string, txHash string, envelope *envelope.Envelope) (status string, at time.Time, err error)

	// Load context envelope by envelope ID
	LoadByID(ctx context.Context, envelopeID string, envelope *envelope.Envelope) (status string, at time.Time, err error)

	// Load context envelopes that have been pending for at least a given duration
	LoadPending(ctx context.Context, duration time.Duration) ([]*envelope.Envelope, error)

	// GetStatus returns envelope status
	GetStatus(ctx context.Context, envelopeID string) (status string, at time.Time, err error)

	// SetStatus set envelope status
	SetStatus(ctx context.Context, envelopeID string, status string) error
}
