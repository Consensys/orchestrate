package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	envelope "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
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
func (s *EnvelopeStore) Store(ctx context.Context, e *envelope.Envelope) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	k := key(e.GetChain().ID().String(), e.GetTx().GetHash().Hex())
	entry, ok := s.byTxHash[k]
	if ok {
		entry.envelope = e
		return entry.Status, entry.last(), nil
	}

	entry = &Entry{
		envelope: e,
		Status:   "stored",
		StoredAt: time.Now(),
	}
	s.byID[e.GetMetadata().GetId()] = entry
	s.byTxHash[k] = entry

	return STORED, entry.StoredAt, nil
}

// LoadByTxHash context envelope by transaction hash
func (s *EnvelopeStore) LoadByTxHash(ctx context.Context, chainID, txHash string, e *envelope.Envelope) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byTxHash[key(chainID, txHash)]
	if !ok {
		return "", time.Time{}, fmt.Errorf("no envelope for chain %q txHash %q", chainID, txHash)
	}
	proto.Merge(e, entry.envelope)

	return entry.Status, entry.last(), nil
}

// LoadByID context envelope by envelope ID
func (s *EnvelopeStore) LoadByID(ctx context.Context, envelopeID string, e *envelope.Envelope) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[envelopeID]
	if !ok {
		return "", time.Time{}, fmt.Errorf("no envelope for ID %q", envelopeID)
	}
	proto.Merge(e, entry.envelope)

	return entry.Status, entry.last(), nil
}

// SetStatus set a context status
func (s *EnvelopeStore) SetStatus(ctx context.Context, envelopeID, status string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[envelopeID]
	if !ok {
		return fmt.Errorf("no envelope for ID %q", envelopeID)
	}

	entry.Status = status
	switch status {
	case STORED:
		entry.StoredAt = time.Now()
	case ERROR:
		entry.ErrorAt = time.Now()
	case PENDING:
		entry.SentAt = time.Now()
	case MINED:
		entry.MinedAt = time.Now()
	}

	return nil
}

// GetStatus return context status and time when status changed
func (s *EnvelopeStore) GetStatus(ctx context.Context, envelopeID string) (status string, at time.Time, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	entry, ok := s.byID[envelopeID]
	if !ok {
		return "", time.Time{}, fmt.Errorf("no envelope for ID %q", envelopeID)
	}

	return entry.Status, entry.last(), nil
}

// LoadPending loads pending envelopes
func (s *EnvelopeStore) LoadPending(ctx context.Context, duration time.Duration) ([]*envelope.Envelope, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	envelopes := []*envelope.Envelope{}
	for _, entry := range s.byTxHash {
		if entry.Status == PENDING && time.Now().Add(-duration).Sub(entry.SentAt) > 0 {
			envelopes = append(envelopes, entry.envelope)
		}
	}

	return envelopes, nil
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
