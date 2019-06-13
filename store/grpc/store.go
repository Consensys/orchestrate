package grpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	envelope "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope-store"
)

// EnvelopeStore store envelopes through
type EnvelopeStore struct {
	client store.StoreClient
}

// NewEnvelopeStore crate a enw envelope store on GRPC
func NewEnvelopeStore(client store.StoreClient) *EnvelopeStore {
	return &EnvelopeStore{client: client}
}

// Store envelope
func (s *EnvelopeStore) Store(ctx context.Context, e *envelope.Envelope) (status string, at time.Time, err error) {
	resp, err := s.client.Store(ctx, &store.StoreRequest{Envelope: e})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	proto.Merge(e, resp.Envelope)

	return resp.Status, at, nil
}

// LoadByTxHash load envelope by TxHash
func (s *EnvelopeStore) LoadByTxHash(ctx context.Context, chainID, txHash string, e *envelope.Envelope) (status string, at time.Time, err error) {
	resp, err := s.client.LoadByTxHash(ctx, &store.TxHashRequest{ChainId: chainID, TxHash: txHash})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	proto.Merge(e, resp.Envelope)

	return resp.Status, at, nil
}

// LoadByID load envelope by TxHash
func (s *EnvelopeStore) LoadByID(ctx context.Context, envelopeID string, e *envelope.Envelope) (status string, at time.Time, err error) {
	resp, err := s.client.LoadByID(ctx, &store.IDRequest{Id: envelopeID})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	return resp.Status, at, nil
}

// SetStatus set envelope status
func (s *EnvelopeStore) SetStatus(ctx context.Context, envelopeID, status string) error {
	_, err := s.client.SetStatus(ctx, &store.SetStatusRequest{Id: envelopeID, Status: status})
	if err != nil {
		return err
	}

	return nil
}

// GetStatus get envelope status
func (s *EnvelopeStore) GetStatus(ctx context.Context, envelopeID string) (status string, at time.Time, err error) {
	resp, err := s.client.GetStatus(ctx, &store.IDRequest{Id: envelopeID})
	if err != nil {
		return "", time.Time{}, err
	}

	at, err = ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		return "", time.Time{}, err
	}

	return resp.Status, at, nil
}

// LoadPending loads pending envelopes
func (s *EnvelopeStore) LoadPending(ctx context.Context, duration time.Duration) ([]*envelope.Envelope, error) {
	resp, err := s.client.LoadPending(ctx, &store.LoadPendingRequest{Duration: duration.Nanoseconds()})

	if err != nil {
		return nil, err
	}

	return resp.Envelopes, nil
}
