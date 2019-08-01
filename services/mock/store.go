package mock

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
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
	byID     map[string]*EnvelopeModel
	byTxHash map[string]*EnvelopeModel
}

// NewEnvelopeStore creates a new mock envelope store
func NewEnvelopeStore() *EnvelopeStore {
	return &EnvelopeStore{
		mux:      &sync.Mutex{},
		byID:     make(map[string]*EnvelopeModel),
		byTxHash: make(map[string]*EnvelopeModel),
	}
}

func key(chainID, txHash string) string {
	return fmt.Sprintf("%v-%v", chainID, txHash)
}

// Store envelope
func (s *EnvelopeStore) Store(ctx context.Context, req *evlpstore.StoreRequest) (*evlpstore.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Prepare new model to be inserted
	model := &EnvelopeModel{
		envelope: req.GetEnvelope(),
		Status:   STORED,
		StoredAt: time.Now(),
	}

	// Store model
	s.byID[req.GetEnvelope().GetMetadata().GetId()] = model

	k := key(req.GetEnvelope().GetChain().ID().String(), req.GetEnvelope().GetTx().GetHash().Hex())
	s.byTxHash[k] = model

	return model.ToStoreResponse()
}

// LoadByTxHash context envelope by transaction hash
func (s *EnvelopeStore) LoadByTxHash(ctx context.Context, req *evlpstore.LoadByTxHashRequest) (*evlpstore.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	model, ok := s.byTxHash[key(req.GetChain().ID().String(), req.GetTxHash().Hex())]
	if !ok {
		return &evlpstore.StoreResponse{},
			errors.NotFoundError("no envelope for chain %q txHash %q", req.GetChain().ID().String(), req.GetTxHash().Hex()).
				SetComponent(component)
	}

	return model.ToStoreResponse()
}

// LoadByID context envelope by envelope ID
func (s *EnvelopeStore) LoadByID(ctx context.Context, req *evlpstore.LoadByIDRequest) (*evlpstore.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	model, ok := s.byID[req.GetId()]
	if !ok {
		return &evlpstore.StoreResponse{},
			errors.NotFoundError("no envelope for ID %q", req.GetId()).
				SetComponent(component)
	}

	return model.ToStoreResponse()
}

// SetStatus set a context status
func (s *EnvelopeStore) SetStatus(ctx context.Context, req *evlpstore.SetStatusRequest) (*evlpstore.StatusResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	model, ok := s.byID[req.GetId()]
	if !ok {
		return &evlpstore.StatusResponse{},
			errors.NotFoundError("no envelope for ID %q", req.GetId()).
				SetComponent(component)
	}

	// Set status
	model.Status = strings.ToLower(req.GetStatus().String())

	// Update time
	switch req.GetStatus() {
	case evlpstore.Status_STORED:
		model.StoredAt = time.Now()
	case evlpstore.Status_ERROR:
		model.ErrorAt = time.Now()
	case evlpstore.Status_PENDING:
		model.SentAt = time.Now()
	case evlpstore.Status_MINED:
		model.MinedAt = time.Now()
	}

	return model.ToStatusResponse()
}

// LoadPending loads pending envelopes
func (s *EnvelopeStore) LoadPending(ctx context.Context, req *evlpstore.LoadPendingRequest) (*evlpstore.LoadPendingResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	resps := []*evlpstore.StoreResponse{}
	for _, model := range s.byID {
		if model.Status == "pending" && time.Now().Add(-utils.PDurationToDuration(req.GetDuration())).Sub(model.SentAt) > 0 {
			resp, err := model.ToStoreResponse()
			if err != nil {
				return &evlpstore.LoadPendingResponse{}, errors.FromError(err).ExtendComponent(component)
			}
			resps = append(resps, resp)
		}
	}

	return &evlpstore.LoadPendingResponse{
		Responses: resps,
	}, nil
}
