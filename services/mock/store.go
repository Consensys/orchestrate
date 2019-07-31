package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

const (
	STORED  = "stored"
	ERROR   = "error"
	PENDING = "pending"
	MINED   = "mined"
)

// EnvelopeStore is a mock in memory version of a envelope store
type EnvelopeStore struct {
	mux      *sync.Mutex
	byID     map[string]*Entry
	byTxHash map[string]*Entry
}

// NewEnvelopeStore creates a new mock envelope store
func NewEnvelopeStore() *EnvelopeStore {
	return &EnvelopeStore{
		mux:      &sync.Mutex{},
		byID:     make(map[string]*Entry),
		byTxHash: make(map[string]*Entry),
	}
}

func key(chainID, txHash string) string {
	return fmt.Sprintf("%v-%v", chainID, txHash)
}

// Store envelope
func (s *EnvelopeStore) Store(ctx context.Context, req *evlpstore.StoreRequest) (*evlpstore.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	resp := &evlpstore.StoreResponse{}

	e := req.GetEnvelope()
	k := key(e.GetChain().ID().String(), e.GetTx().GetHash().Hex())
	entry, ok := s.byTxHash[k]
	if ok {
		// Prepare response with already stored envelope
		resp.Envelope = entry.envelope

		// Update envelope
		entry.envelope = e
	} else {
		// Create entry
		entry = &Entry{
			envelope: e,
			Status:   STORED,
			StoredAt: time.Now(),
		}

		// Store entry
		s.byID[e.GetMetadata().GetId()] = entry
		s.byTxHash[k] = entry
	}

	// Set response attributes
	timestamp, err := ptypes.TimestampProto(entry.last())
	if err != nil {
		panic(err)
	}
	resp.LastUpdated = timestamp
	resp.Status = entry.Status

	return resp, nil
}

// LoadByTxHash context envelope by transaction hash
func (s *EnvelopeStore) LoadByTxHash(ctx context.Context, req *evlpstore.TxHashRequest) (*evlpstore.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byTxHash[key(req.GetChainId().ID().String(), req.GetTxHash())]
	if !ok {
		return &evlpstore.StoreResponse{},
			errors.NotFoundError("no envelope for chain %q txHash %q", req.GetChainId().ID().String(), req.GetTxHash()).
				SetComponent(component)
	}

	timestamp, err := ptypes.TimestampProto(entry.last())
	if err != nil {
		panic(err)
	}

	return &evlpstore.StoreResponse{
		Envelope:    entry.envelope,
		Status:      entry.Status,
		LastUpdated: timestamp,
	}, nil
}

// LoadByID context envelope by envelope ID
func (s *EnvelopeStore) LoadByID(ctx context.Context, req *evlpstore.IDRequest) (*evlpstore.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[req.GetId()]
	if !ok {
		return &evlpstore.StoreResponse{},
			errors.NotFoundError("no envelope for ID %q", req.GetId()).
				SetComponent(component)
	}

	timestamp, err := ptypes.TimestampProto(entry.last())
	if err != nil {
		panic(err)
	}

	return &evlpstore.StoreResponse{
		Envelope:    entry.envelope,
		Status:      entry.Status,
		LastUpdated: timestamp,
	}, nil
}

// SetStatus set a context status
func (s *EnvelopeStore) SetStatus(ctx context.Context, req *evlpstore.SetStatusRequest) (*evlpstore.SetStatusResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[req.GetId()]
	if !ok {
		return &evlpstore.SetStatusResponse{},
			errors.NotFoundError("no envelope for ID %q", req.GetId()).
				SetComponent(component)
	}

	entry.Status = req.GetStatus()
	switch entry.Status {
	case STORED:
		entry.StoredAt = time.Now()
	case ERROR:
		entry.ErrorAt = time.Now()
	case PENDING:
		entry.SentAt = time.Now()
	case MINED:
		entry.MinedAt = time.Now()
	}

	return &evlpstore.SetStatusResponse{}, nil
}

// GetStatus return context status and time when status changed
func (s *EnvelopeStore) GetStatus(ctx context.Context, req *evlpstore.IDRequest) (*evlpstore.StoreResponse, error) {
	return s.LoadByID(ctx, req)
}

// LoadPending loads pending envelopes
func (s *EnvelopeStore) LoadPending(ctx context.Context, req *evlpstore.LoadPendingRequest) (*evlpstore.LoadPendingResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	envelopes := []*envelope.Envelope{}
	for _, entry := range s.byTxHash {
		if entry.Status == PENDING && time.Now().Add(-time.Duration(req.GetDuration())).Sub(entry.SentAt) > 0 {
			envelopes = append(envelopes, entry.envelope)
		}
	}

	return &evlpstore.LoadPendingResponse{
		Envelopes: envelopes,
	}, nil
}

// Entry is a entry into mock envelope Store
type Entry struct {
	// envelope
	envelope *envelope.Envelope

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
