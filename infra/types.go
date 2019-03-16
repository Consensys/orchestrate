package infra

import (
	"context"
	"time"

	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// TraceStore is an interface for context store
type TraceStore interface {
	// Store context trace
	// For a trace with same TxHash that was already store, it updates the trace and return status
	Store(ctx context.Context, trace *trace.Trace) (status string, at time.Time, err error)

	// Load context trace by txHash
	LoadByTxHash(ctx context.Context, chainID string, txHash string, trace *trace.Trace) (status string, at time.Time, err error)

	// Load context trace by trace ID
	LoadByTraceID(ctx context.Context, traceID string, trace *trace.Trace) (status string, at time.Time, err error)

	// Load context traces that have been pending for at least a given duration
	LoadPendingTraces(ctx context.Context, duration time.Duration) ([]*trace.Trace, error)

	// GetStatus returns trace status
	GeStatus(ctx context.Context, traceID string) (status string, at time.Time, err error)

	// SetStatus set trace status
	SeStatus(ctx context.Context, traceID string, status string) error
}
