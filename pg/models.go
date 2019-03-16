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

// TraceModel represent elements in `traces` table
type TraceModel struct {
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

	model := &TraceModel{
		ChainID: trace.GetChain().GetId(),
		TxHash:  trace.GetTx().Hash,
		TraceID: trace.GetMetadata().GetId(),
		Trace:   bytes,
	}

	_, err = s.db.ModelContext(ctx, model).
		OnConflict("ON CONSTRAINT uni_tx DO UPDATE").
		Set("trace = ?trace").
		Returning("*").
		Insert()
	if err != nil {
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

// LoadByTxHash context trace by transaction hash
func (s *traceStore) LoadByTxHash(ctx context.Context, chainID string, txHash string, trace *trace.Trace) (status string, at time.Time, err error) {
	model := &TraceModel{
		ChainID: chainID,
		TxHash:  txHash,
	}
	return s.Load(ctx, model, trace)
}

// LoadByTraceID context trace by trace ID
func (s *traceStore) LoadByTraceID(ctx context.Context, traceID string, trace *trace.Trace) (status string, at time.Time, err error) {
	model := &TraceModel{
		TraceID: traceID,
	}
	return s.Load(ctx, model, trace)
}

func (s *traceStore) Load(ctx context.Context, model *TraceModel, trace *trace.Trace) (status string, at time.Time, err error) {
	err = s.db.ModelContext(ctx, model).Select()
	if err != nil {
		return "", time.Time{}, err
	}

	err = proto.UnmarshalMerge(model.Trace, trace)
	if err != nil {
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

// SetStatus set a context status
func (s *traceStore) SetStatus(ctx context.Context, traceID string, status string) error {
	// Define model
	model := &TraceModel{
		TraceID: traceID,
		Status:  status,
	}

	// Update status value
	_, err := s.db.ModelContext(ctx, model).
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
	model := &TraceModel{
		TraceID: traceID,
	}

	err = s.db.ModelContext(ctx, model).Select()
	if err != nil {
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

func (model *TraceModel) last() time.Time {
	switch model.Status {
	case "stored":
		return model.StoredAt
	case "error":
		return model.ErrorAt
	case "pending":
		return model.SentAt
	case "mined":
		return model.MinedAt
	}
	return time.Time{}
}
