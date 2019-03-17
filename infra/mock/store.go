package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// TraceStore is a mock in memory version of a trace store
type TraceStore struct {
	mux      *sync.Mutex
	byID     map[string]*Entry
	byTxHash map[string]*Entry
}

// NewTraceStore creates a new mock trace store
func NewTraceStore() *TraceStore {
	return &TraceStore{
		mux:      &sync.Mutex{},
		byID:     make(map[string]*Entry),
		byTxHash: make(map[string]*Entry),
	}
}

func key(chainID string, TxHash string) string {
	return fmt.Sprintf("%v-%v", chainID, TxHash)
}

// Store trace
func (s *TraceStore) Store(ctx context.Context, trace *trace.Trace) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	k := key(trace.GetChain().GetId(), trace.GetTx().GetHash())
	entry, ok := s.byTxHash[k]
	if ok {
		entry.Trace = trace
		return entry.Status, entry.last(), nil
	}

	entry = &Entry{
		Trace:    trace,
		Status:   "stored",
		StoredAt: time.Now(),
	}
	s.byID[trace.GetMetadata().GetId()] = entry
	s.byTxHash[k] = entry

	return "stored", entry.StoredAt, nil
}

// LoadByTxHash context trace by transaction hash
func (s *TraceStore) LoadByTxHash(ctx context.Context, chainID string, txHash string, trace *trace.Trace) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byTxHash[key(chainID, txHash)]
	if !ok {
		return "", time.Time{}, fmt.Errorf("No trace for chain %q txHash %q", chainID, txHash)
	}
	proto.Merge(trace, entry.Trace)

	return entry.Status, entry.last(), nil
}

// LoadByTraceID context trace by trace ID
func (s *TraceStore) LoadByTraceID(ctx context.Context, traceID string, trace *trace.Trace) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[traceID]
	if !ok {
		return "", time.Time{}, fmt.Errorf("No trace for ID %q", traceID)
	}
	proto.Merge(trace, entry.Trace)

	return entry.Status, entry.last(), nil
}

// SetStatus set a context status
func (s *TraceStore) SetStatus(ctx context.Context, traceID string, status string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[traceID]
	if !ok {
		return fmt.Errorf("No trace for ID %q", traceID)
	}

	entry.Status = status
	switch status {
	case "stored":
		entry.StoredAt = time.Now()
	case "error":
		entry.ErrorAt = time.Now()
	case "pending":
		entry.SentAt = time.Now()
	case "mined":
		entry.MinedAt = time.Now()
	}

	return nil
}

// GetStatus return context status and time when status changed
func (s *TraceStore) GetStatus(ctx context.Context, traceID string) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[traceID]
	if !ok {
		return "", time.Time{}, fmt.Errorf("No trace for ID %q", traceID)
	}

	return entry.Status, entry.last(), nil
}

// LoadPendingTraces loads pending traces
func (s *TraceStore) LoadPendingTraces(ctx context.Context, duration time.Duration) ([]*trace.Trace, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	traces := []*trace.Trace{}
	for _, entry := range s.byTxHash {
		if entry.Status == "pending" && time.Now().Add(-duration).Sub(entry.SentAt) > 0 {
			traces = append(traces, entry.Trace)
		}
	}

	return traces, nil
}

// Entry is a entry into mock Trace Store
type Entry struct {
	// Trace
	Trace *trace.Trace

	// Status
	Status   string
	StoredAt time.Time
	ErrorAt  time.Time
	SentAt   time.Time
	MinedAt  time.Time
}

func (entry *Entry) last() time.Time {
	switch entry.Status {
	case "stored":
		return entry.StoredAt
	case "error":
		return entry.ErrorAt
	case "pending":
		return entry.SentAt
	case "mined":
		return entry.MinedAt
	}
	return time.Time{}
}
