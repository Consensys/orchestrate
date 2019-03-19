package grpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/context-store"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// TraceStore store traces through
type TraceStore struct {
	client store.StoreClient
}

// NewTraceStore crate a enw trace store on GRPC
func NewTraceStore(client store.StoreClient) *TraceStore {
	return &TraceStore{client: client}
}

// Store trace
func (s *TraceStore) Store(ctx context.Context, trace *trace.Trace) (status string, at time.Time, err error) {
	resp, err := s.client.Store(ctx, &store.StoreRequest{Trace: trace})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	proto.Merge(trace, resp.Trace)

	return resp.Status, at, nil
}

// LoadByTxHash load trace by TxHash
func (s *TraceStore) LoadByTxHash(ctx context.Context, chainID string, txHash string, trace *trace.Trace) (status string, at time.Time, err error) {
	resp, err := s.client.LoadByTxHash(ctx, &store.TxHashRequest{ChainId: chainID, TxHash: txHash})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	proto.Merge(trace, resp.Trace)

	return resp.Status, at, nil
}

// LoadByTraceID load trace by TxHash
func (s *TraceStore) LoadByTraceID(ctx context.Context, traceID string, trace *trace.Trace) (status string, at time.Time, err error) {
	resp, err := s.client.LoadByTraceID(ctx, &store.TraceIDRequest{TraceId: traceID})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	return resp.Status, at, nil
}

// SetStatus set trace status
func (s *TraceStore) SetStatus(ctx context.Context, traceID string, status string) error {
	_, err := s.client.SetStatus(ctx, &store.SetStatusRequest{TraceId: traceID, Status: status})
	if err != nil {
		return err
	}

	return nil
}

// GetStatus get trace status
func (s *TraceStore) GetStatus(ctx context.Context, traceID string) (status string, at time.Time, err error) {
	resp, err := s.client.GetStatus(ctx, &store.TraceIDRequest{TraceId: traceID})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	return resp.Status, at, nil
}

// LoadPendingTraces loads pending traces
func (s *TraceStore) LoadPendingTraces(ctx context.Context, duration time.Duration) ([]*trace.Trace, error) {
	resp, err := s.client.LoadPendingTraces(ctx, &store.PendingTracesRequest{Duration: duration.Nanoseconds()})

	if err != nil {
		return nil, err
	}

	return resp.Traces, nil
}
