package pg

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/golang/protobuf/proto"
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
	LoadPendingTraces(duration time.Duration) ([]*trace.Trace, error)

	// GetStatus returns trace status
	GeStatus(ctx context.Context, traceID string) (status string, at time.Time, err error)

	// SetStatus set trace status
	SeStatus(ctx context.Context, traceID string, status string) error
}

// Trace represent elements in `traces` table
type Trace struct {
	tableName struct{} `sql:"traces"`

	// ID technical identifier
	ID int32

	// Context Identifier
	ChainID string
	TxHash  string
	TraceID string

	// Trace
	Trace []byte

	// Status
	Status   string
	StoredAt time.Time
	ErrorAt  time.Time
	SentAt   time.Time
	MinedAt  time.Time
}

// traceStore is a context store based on PostgreSQL
type traceStore struct {
	db *pg.DB
}

// Store context trace
func (s *traceStore) Store(ctx context.Context, trace *trace.Trace) (status string, at time.Time, err error) {
	bytes, err := proto.Marshal(trace)
	if err != nil {
		return "", time.Time{}, err
	}

	t := &Trace{
		ChainID: trace.GetChain().GetId(),
		TxHash:  trace.GetTx().Hash,
		TraceID: trace.GetMetadata().GetId(),
		Trace:   bytes,
	}

	_, err = s.db.ModelContext(ctx, t).Returning("*").Insert()
	if err != nil {
		return "", time.Time{}, err
	}

	return t.Status, last(t), nil
}

func last(t *Trace) time.Time {
	switch t.Status {
	case "stored":
		return t.StoredAt
	case "error":
		return t.ErrorAt
	case "pending":
		return t.SentAt
	case "mined":
		return t.MinedAt
	}
	return time.Time{}
}

// LoadByTxHash context trace by transaction hash
func (s *traceStore) LoadByTxHash(ctx context.Context, chainID string, txHash string, trace *trace.Trace) (status string, at time.Time, err error) {
	t := &Trace{
		TxHash:  txHash,
		ChainID: chainID,
	}

	err = s.db.ModelContext(ctx, t).Select()
	if err != nil {
		return "", time.Time{}, err
	}

	err = proto.UnmarshalMerge(t.Trace, trace)
	if err != nil {
		return "", time.Time{}, err
	}

	return t.Status, last(t), nil
}

// LoadByTraceID context trace by trace ID
func (s *traceStore) LoadByTraceID(ctx context.Context, traceID string, trace *trace.Trace) (status string, at time.Time, err error) {
	t := &Trace{
		TraceID: traceID,
	}

	err = s.db.ModelContext(ctx, t).Select()
	if err != nil {
		return "", time.Time{}, err
	}

	err = proto.UnmarshalMerge(t.Trace, trace)
	if err != nil {
		return "", time.Time{}, err
	}

	return t.Status, last(t), nil
}

// SetStatus set a context status
func (s *traceStore) SetStatus(ctx context.Context, traceID string, status string) error {
	t := &Trace{
		TraceID: traceID,
		Status:  status,
	}

	_, err := s.db.ModelContext(ctx, t).
		Set("status = ?status").
		Where("trace_id = ?trace_id").
		Update()
	if err != nil {
		return err
	}

	return nil
}

// GetStatus return context status and time when status changed
func (s *traceStore) GetStatus(ctx context.Context, traceID string) (status string, at time.Time, err error) {
	t := &Trace{
		TraceID: traceID,
	}

	err = s.db.ModelContext(ctx, t).Select()
	if err != nil {
		return "", time.Time{}, err
	}

	return t.Status, last(t), nil
}
