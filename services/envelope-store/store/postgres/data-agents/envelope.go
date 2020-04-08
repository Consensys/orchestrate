package dataagents

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

// PGEnvelope is a tag data agent
type PGEnvelopeAgent struct {
	db *pg.DB
}

// NewPGEnvelope creates a new PGEnvelope
func NewPGEnvelope(db *pg.DB) *PGEnvelopeAgent {
	return &PGEnvelopeAgent{db: db}
}

func (ag *PGEnvelopeAgent) InsertDoUpdateOnEnvelopeIDKey(ctx context.Context, obj *models.EnvelopeModel) error {
	// Execute ORM query
	// If uniqueness constraint is broken then it update the former value
	_, err := ag.db.ModelContext(ctx, obj).
		OnConflict("ON CONSTRAINT envelopes_envelope_id_key DO UPDATE").
		Set("envelope = ?envelope").
		Set("chain_id = ?chain_id").
		Set("tx_hash = ?tx_hash").
		Returning("*").
		Insert()

	if err != nil {
		return err
	}

	return nil
}

func (ag *PGEnvelopeAgent) InsertDoUpdateOnUniTx(ctx context.Context, obj *models.EnvelopeModel) error {
	// Possibly we got an error due to unique contraint on tx,chain_id so we try again
	_, err := ag.db.ModelContext(ctx, obj).
		OnConflict("ON CONSTRAINT uni_tx DO UPDATE").
		Set("envelope = ?envelope").
		Set("envelope_id = ?envelope_id").
		Returning("*").
		Insert()

	if err != nil {
		return err
	}

	return nil
}

func (ag *PGEnvelopeAgent) FindByFieldSet(ctx context.Context, fields map[string]string) (models.EnvelopeModel, error) {
	model := models.EnvelopeModel{}
	q := ag.db.ModelContext(ctx, model)
	for key, val := range fields {
		q.Where(fmt.Sprintf("%s = ?", key), val)
	}

	err := q.Select()
	return model, err
}

func (ag *PGEnvelopeAgent) FindPending(ctx context.Context, sentBeforeAt time.Time) ([]*models.EnvelopeModel, error) {
	var envelopes []*models.EnvelopeModel
	err := ag.db.ModelContext(ctx, &envelopes).
		Where("status = 'pending'").
		Where("sent_at < ?", sentBeforeAt).
		Select()

	return envelopes, err
}

func (ag *PGEnvelopeAgent) UpdateStatus(ctx context.Context, envelope *models.EnvelopeModel) error {
	_, err := ag.db.ModelContext(ctx, envelope).
		Set("status = ?status").
		Where("envelope_id = ?", envelope.EnvelopeID).
		Where("tenant_id = ?", envelope.TenantID).
		Returning("*").
		Update()

	return err
}
