package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

const (
	STORED  = "stored"
	ERROR   = "error"
	PENDING = "pending"
	MINED   = "mined"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context) (*Memory, error) {
	return New(), nil
}

// Memory is a in memory version of a envelope store
type Memory struct {
	mux      *sync.Mutex
	byID     map[string]*EnvelopeModel
	byTxHash map[string]*EnvelopeModel
}

// New creates a new in memory envelope store
func New() *Memory {
	return &Memory{
		mux:      &sync.Mutex{},
		byID:     make(map[string]*EnvelopeModel),
		byTxHash: make(map[string]*EnvelopeModel),
	}
}

func key(chainID, txHash string) string {
	return fmt.Sprintf("%v-%v", chainID, txHash)
}

// Store envelope
func (s *Memory) Store(ctx context.Context, req *svc.StoreRequest) (*svc.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Prepare new model to be inserted
	model := &EnvelopeModel{
		envelope: req.GetEnvelope(),
		Status:   STORED,
		StoredAt: time.Now(),
	}

	// Store model
	s.byID[req.GetEnvelope().GetID()] = model

	k := key(req.GetEnvelope().GetChainID(), req.GetEnvelope().GetTxHash())
	s.byTxHash[k] = model

	return model.ToStoreResponse()
}

// LoadByTxHash context envelope by transaction hash
func (s *Memory) LoadByTxHash(ctx context.Context, req *svc.LoadByTxHashRequest) (*svc.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	model, ok := s.byTxHash[key(req.GetChainId(), req.GetTxHash())]
	if !ok {
		return &svc.StoreResponse{}, errors.NotFoundError("no envelope for chain %q txHash %q", req.GetChainId(), req.GetTxHash())
	}

	return model.ToStoreResponse()
}

// LoadByID context envelope by envelope UUID
func (s *Memory) LoadByID(ctx context.Context, req *svc.LoadByIDRequest) (*svc.StoreResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	model, ok := s.byID[req.GetId()]
	if !ok {
		return &svc.StoreResponse{},
			errors.NotFoundError("no envelope for UUID %q", req.GetId())

	}

	return model.ToStoreResponse()
}

// SetStatus set a context status
func (s *Memory) SetStatus(ctx context.Context, req *svc.SetStatusRequest) (*svc.StatusResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	model, ok := s.byID[req.GetId()]
	if !ok {
		return &svc.StatusResponse{},
			errors.NotFoundError("no envelope for UUID %q", req.GetId())

	}

	// Set status
	model.Status = strings.ToLower(req.GetStatus().String())

	// Update time
	switch req.GetStatus() {
	case svc.Status_STORED:
		model.StoredAt = time.Now()
	case svc.Status_ERROR:
		model.ErrorAt = time.Now()
	case svc.Status_PENDING:
		model.SentAt = time.Now()
	case svc.Status_MINED:
		model.MinedAt = time.Now()
	}

	return model.ToStatusResponse()
}

// LoadPending loads pending envelopes
func (s *Memory) LoadPending(ctx context.Context, req *svc.LoadPendingRequest) (*svc.LoadPendingResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var resps []*svc.StoreResponse
	for _, model := range s.byID {
		if model.Status == "pending" && time.Now().Add(-utils.PDurationToDuration(req.GetDuration())).Sub(model.SentAt) > 0 {
			resp, err := model.ToStoreResponse()
			if err != nil {
				return &svc.LoadPendingResponse{}, errors.FromError(err)
			}
			resps = append(resps, resp)
		}
	}

	return &svc.LoadPendingResponse{
		Responses: resps,
	}, nil
}
