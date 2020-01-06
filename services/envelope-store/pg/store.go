package pg

import (
	"context"
	"strings"
	"time"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
)

// EnvelopeStore is a context store based on PostgreSQL
type EnvelopeStore struct {
	db *pg.DB
}

// NewEnvelopeStore creates a new envelope store
func NewEnvelopeStore(db *pg.DB) *EnvelopeStore {
	return &EnvelopeStore{db: db}
}

// NewEnvelopeStoreFromPGOptions creates a new pg envelope store
func NewEnvelopeStoreFromPGOptions(opts *pg.Options) *EnvelopeStore {
	return NewEnvelopeStore(pg.Connect(opts))
}

// Store context envelope
func (s *EnvelopeStore) Store(ctx context.Context, req *evlpstore.StoreRequest) (*evlpstore.StoreResponse, error) {
	// create model from envelope
	model, err := FromEnvelope(req.GetEnvelope())
	if err != nil {
		return &evlpstore.StoreResponse{}, errors.FromError(err).SetComponent(component)
	}

	// Execute ORM query
	// If uniqueness constraint is broken then it update the former value
	_, err = s.db.ModelContext(ctx, model).
		OnConflict("ON CONSTRAINT envelopes_envelope_id_key DO UPDATE").
		Set("envelope = ?envelope").
		Set("chain_id = ?chain_id").
		Set("tx_hash = ?tx_hash").
		Returning("*").
		Insert()
	if err != nil {
		// Possibly we got an error due to unique contraint on tx,chain_id so we try again
		_, err = s.db.ModelContext(ctx, model).
			OnConflict("ON CONSTRAINT uni_tx DO UPDATE").
			Set("envelope = ?envelope").
			Set("envelope_id = ?envelope_id").
			Returning("*").
			Insert()
		if err != nil {
			log.WithError(err).Error("Could not store")
			return &evlpstore.StoreResponse{}, errors.StorageError("%v", err).ExtendComponent(component)
		}
	}

	return model.ToStoreResponse()
}

// LoadByTxHash load envelope by transaction hash
func (s *EnvelopeStore) LoadByTxHash(ctx context.Context, req *evlpstore.LoadByTxHashRequest) (*evlpstore.StoreResponse, error) { //nolint:interfacer // reason
	model := &EnvelopeModel{
		ChainID: req.GetChain().ID().String(),
		TxHash:  req.GetTxHash().Hex(),
	}

	err := s.db.ModelContext(ctx, model).
		Where("chain_id = ?", model.ChainID).
		Where("tx_hash = ?", model.TxHash).
		Select()
	if err != nil {
		return &evlpstore.StoreResponse{}, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	return model.ToStoreResponse()
}

// LoadByID context envelope by envelope ID
func (s *EnvelopeStore) LoadByID(ctx context.Context, req *evlpstore.LoadByIDRequest) (*evlpstore.StoreResponse, error) { //nolint:interfacer // reason
	model := &EnvelopeModel{
		EnvelopeID: req.GetId(),
	}

	err := s.db.ModelContext(ctx, model).
		Where("envelope_id = ?", model.EnvelopeID).
		Select()
	if err != nil {
		return &evlpstore.StoreResponse{}, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	return model.ToStoreResponse()
}

// SetStatus set a context status
func (s *EnvelopeStore) SetStatus(ctx context.Context, req *evlpstore.SetStatusRequest) (*evlpstore.StatusResponse, error) {
	// Define model
	model := &EnvelopeModel{
		EnvelopeID: req.GetId(),
		Status:     strings.ToLower(req.GetStatus().String()),
	}

	// Update status value
	_, err := s.db.ModelContext(ctx, model).
		Set("status = ?status").
		Where("envelope_id = ?envelope_id").
		Returning("*").
		Update()
	if err != nil {
		return &evlpstore.StatusResponse{}, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	return model.ToStatusResponse()
}

// LoadPending loads pending envelopes
func (s *EnvelopeStore) LoadPending(ctx context.Context, req *evlpstore.LoadPendingRequest) (*evlpstore.LoadPendingResponse, error) {
	var models []*EnvelopeModel

	err := s.db.ModelContext(ctx, &models).
		Where("status = 'pending'").
		Where("sent_at < ?", time.Now().Add(-utils.PDurationToDuration(req.GetDuration()))).
		Select()
	if err != nil {
		return nil, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	var resps []*evlpstore.StoreResponse
	for _, model := range models {
		resp, err := model.ToStoreResponse()
		if err != nil {
			return &evlpstore.LoadPendingResponse{}, errors.FromError(err).ExtendComponent(component)
		}
		resps = append(resps, resp)
	}

	return &evlpstore.LoadPendingResponse{
		Responses: resps,
	}, nil
}
