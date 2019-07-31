package pg

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/golang/protobuf/ptypes"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// EnvelopeModel represent elements in `envelopes` table
type EnvelopeModel struct {
	tableName struct{} `sql:"envelopes"` //nolint:unused,structcheck

	// ID technical identifier
	ID int32

	// Context Identifier
	ChainID    string
	TxHash     string
	EnvelopeID string

	// Envelope
	Envelope []byte

	// Status
	Status   string
	StoredAt time.Time
	ErrorAt  time.Time
	SentAt   time.Time
	MinedAt  time.Time
}

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

func ToStoreResponse(model *EnvelopeModel) (*evlpstore.StoreResponse, error) {
	// Prepare response
	timestamp, err := ptypes.TimestampProto(model.last())
	if err != nil {
		panic(err)
	}

	resp := &evlpstore.StoreResponse{
		Envelope:    &envelope.Envelope{},
		LastUpdated: timestamp,
		Status:      model.Status,
	}

	err = encoding.Unmarshal(model.Envelope, resp.GetEnvelope())
	if err != nil {
		return resp, errors.FromError(err).SetComponent(component)
	}

	return resp, err
}

// Store context envelope
func (s *EnvelopeStore) Store(ctx context.Context, req *evlpstore.StoreRequest) (*evlpstore.StoreResponse, error) {
	e := req.GetEnvelope()
	bytes, err := encoding.Marshal(e)
	if err != nil {
		return &evlpstore.StoreResponse{}, errors.FromError(err).SetComponent(component)
	}

	// Execute ORM query
	// If unicity contraint is broken then it update the former value
	model := &EnvelopeModel{
		ChainID:    e.GetChain().ID().String(),
		TxHash:     e.GetTx().GetHash().Hex(),
		EnvelopeID: e.GetMetadata().GetId(),
		Envelope:   bytes,
	}
	_, err = s.db.ModelContext(ctx, model).
		OnConflict("ON CONSTRAINT uni_tx DO UPDATE").
		Set("envelope = ?envelope").
		Returning("*").
		Insert()
	if err != nil {
		return &evlpstore.StoreResponse{}, errors.ConstraintViolatedError("envelope already stored").ExtendComponent(component)
	}

	return ToStoreResponse(model)
}

// LoadByTxHash load envelope by transaction hash
func (s *EnvelopeStore) LoadByTxHash(ctx context.Context, req *evlpstore.TxHashRequest) (*evlpstore.StoreResponse, error) { //nolint:interfacer
	model := &EnvelopeModel{
		ChainID: req.GetChainId().ID().String(),
		TxHash:  req.GetTxHash(),
	}

	err := s.db.ModelContext(ctx, model).
		Where("chain_id = ?", model.ChainID).
		Where("tx_hash = ?", model.TxHash).
		Select()
	if err != nil {
		return &evlpstore.StoreResponse{}, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	return ToStoreResponse(model)
}

// LoadByID context envelope by envelope ID
func (s *EnvelopeStore) LoadByID(ctx context.Context, req *evlpstore.IDRequest) (*evlpstore.StoreResponse, error) { //nolint:interfacer
	model := &EnvelopeModel{
		EnvelopeID: req.GetId(),
	}

	err := s.db.ModelContext(ctx, model).
		Where("envelope_id = ?", model.EnvelopeID).
		Select()
	if err != nil {
		return &evlpstore.StoreResponse{}, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	return ToStoreResponse(model)
}

// SetStatus set a context status
func (s *EnvelopeStore) SetStatus(ctx context.Context, req *evlpstore.SetStatusRequest) (*evlpstore.SetStatusResponse, error) {
	// Define model
	model := &EnvelopeModel{
		EnvelopeID: req.GetId(),
		Status:     req.GetStatus(),
	}

	// Update status value
	_, err := s.db.ModelContext(ctx, model).
		Set("status = ?status").
		Where("envelope_id = ?envelope_id").
		Returning("*").
		Update()
	if err != nil {
		return &evlpstore.SetStatusResponse{}, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	return &evlpstore.SetStatusResponse{}, nil
}

// GetStatus return context status and time when status changed
func (s *EnvelopeStore) GetStatus(ctx context.Context, req *evlpstore.IDRequest) (*evlpstore.StoreResponse, error) {
	return s.LoadByID(ctx, req)
}

// LoadPending loads pending envelopes
func (s *EnvelopeStore) LoadPending(ctx context.Context, req *evlpstore.LoadPendingRequest) (*evlpstore.LoadPendingResponse, error) {
	models := []*EnvelopeModel{}
	err := s.db.ModelContext(ctx, &models).
		Where("status = 'pending'").
		Where("sent_at < ?", time.Now().Add(-time.Duration(req.GetDuration()))).
		Select()

	if err != nil {
		return nil, errors.NotFoundError("envelope not found").ExtendComponent(component)
	}

	envelopes := []*envelope.Envelope{}
	for _, model := range models {
		t := &envelope.Envelope{}
		err := encoding.Unmarshal(model.Envelope, t)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}
		envelopes = append(envelopes, t)
	}

	return &evlpstore.LoadPendingResponse{
		Envelopes: envelopes,
	}, nil
}

func (model *EnvelopeModel) last() time.Time {
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
