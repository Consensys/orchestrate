package pg

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/golang/protobuf/proto"
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

// Store context envelope
func (s *EnvelopeStore) Store(ctx context.Context, e *envelope.Envelope) (status string, at time.Time, err error) {
	bytes, err := proto.Marshal(e)
	if err != nil {
		return "", time.Time{}, err
	}

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
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

// LoadByTxHash context envelope by transaction hash
func (s *EnvelopeStore) LoadByTxHash(ctx context.Context, chainID, txHash string, e *envelope.Envelope) (status string, at time.Time, err error) { //nolint:interfacer
	model := &EnvelopeModel{
		ChainID: chainID,
		TxHash:  txHash,
	}

	err = s.db.ModelContext(ctx, model).
		Where("chain_id = ?", model.ChainID).
		Where("tx_hash = ?", model.TxHash).
		Select()

	if err != nil {
		return "", time.Time{}, err
	}

	err = proto.UnmarshalMerge(model.Envelope, e)
	if err != nil {
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

// LoadByID context envelope by envelope ID
func (s *EnvelopeStore) LoadByID(ctx context.Context, envelopeID string, e *envelope.Envelope) (status string, at time.Time, err error) { //nolint:interfacer
	model := &EnvelopeModel{
		EnvelopeID: envelopeID,
	}

	err = s.db.ModelContext(ctx, model).
		Where("envelope_id = ?", model.EnvelopeID).
		Select()

	if err != nil {
		return "", time.Time{}, err
	}

	err = proto.UnmarshalMerge(model.Envelope, e)
	if err != nil {
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

// SetStatus set a context status
func (s *EnvelopeStore) SetStatus(ctx context.Context, envelopeID, status string) error {
	// Define model
	model := &EnvelopeModel{
		EnvelopeID: envelopeID,
		Status:     status,
	}

	// Update status value
	_, err := s.db.ModelContext(ctx, model).
		Set("status = ?status").
		Where("envelope_id = ?envelope_id").
		Returning("*").
		Update()
	if err != nil {
		return err
	}

	return nil
}

// GetStatus return context status and time when status changed
func (s *EnvelopeStore) GetStatus(ctx context.Context, envelopeID string) (status string, at time.Time, err error) {
	model := &EnvelopeModel{
		EnvelopeID: envelopeID,
	}

	err = s.db.ModelContext(ctx, model).Select()
	if err != nil {
		return "", time.Time{}, err
	}

	return model.Status, model.last(), nil
}

// LoadPending loads pending envelopes
func (s *EnvelopeStore) LoadPending(ctx context.Context, duration time.Duration) ([]*envelope.Envelope, error) {
	models := []*EnvelopeModel{}
	err := s.db.ModelContext(ctx, &models).
		Where("status = 'pending'").
		Where("sent_at < ?", time.Now().Add(-duration)).
		Select()

	if err != nil {
		return nil, err
	}

	envelopes := []*envelope.Envelope{}
	for _, model := range models {
		t := &envelope.Envelope{}
		err := proto.Unmarshal(model.Envelope, t)
		if err != nil {
			return nil, err
		}
		envelopes = append(envelopes, t)
	}

	return envelopes, nil
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
