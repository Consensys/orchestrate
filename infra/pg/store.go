package pg

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/golang/protobuf/proto"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

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

// TraceStore is a context store based on PostgreSQL
type TraceStore struct {
	db *pg.DB
}

// NewTraceStore creates a new trace store
func NewTraceStore(db *pg.DB) *TraceStore {
	return &TraceStore{db: db}
}

// NewTraceStoreFromPGOptions creates a new pg trace store
func NewTraceStoreFromPGOptions(opts *pg.Options) *TraceStore {
	return NewTraceStore(pg.Connect(opts))
}

// Store context trace
func (s *TraceStore) Store(ctx context.Context, trace *trace.Trace) (status string, at time.Time, err error) {
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
func (s *TraceStore) LoadByTxHash(ctx context.Context, chainID string, txHash string, trace *trace.Trace) (status string, at time.Time, err error) {
	model := &TraceModel{
		ChainID: chainID,
		TxHash:  txHash,
	}
	return s.load(ctx, model, trace)
}

// LoadByTraceID context trace by trace ID
func (s *TraceStore) LoadByTraceID(ctx context.Context, traceID string, trace *trace.Trace) (status string, at time.Time, err error) {
	model := &TraceModel{
		TraceID: traceID,
	}
	return s.load(ctx, model, trace)
}

func (s *TraceStore) load(ctx context.Context, model *TraceModel, trace *trace.Trace) (status string, at time.Time, err error) {
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
func (s *TraceStore) SetStatus(ctx context.Context, traceID string, status string) error {
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
func (s *TraceStore) GetStatus(ctx context.Context, traceID string) (status string, at time.Time, err error) {
	model := &TraceModel{
		TraceID: traceID,
	}

	err = s.db.ModelContext(ctx, model).Select()
	if err != nil {
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

// LoadPendingTraces loads pending traces
func (s *TraceStore) LoadPendingTraces(ctx context.Context, duration time.Duration) ([]*trace.Trace, error) {
	models := []*TraceModel{}
	err := s.db.ModelContext(ctx, &models).
		Where("status = 'pending'").
		Where("sent_at < ?", time.Now().Add(-duration)).
		Select()

	if err != nil {
		return nil, err
	}

	traces := []*trace.Trace{}
	for _, model := range models {
		t := &trace.Trace{}
		err := proto.Unmarshal(model.Trace, t)
		if err != nil {
			return nil, err
		}
		traces = append(traces, t)
	}

	return traces, nil
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
