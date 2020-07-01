package store

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

type DataAgents struct {
	Envelope EnvelopeAgent
}

// Interfaces data agents
type EnvelopeAgent interface {
	InsertDoUpdateOnEnvelopeIDKey(ctx context.Context, obj *models.EnvelopeModel) error
	InsertDoUpdateOnUniTx(ctx context.Context, obj *models.EnvelopeModel) error
	FindByFieldSet(ctx context.Context, fields map[string]string, tenantIDs []string) (*models.EnvelopeModel, error)
	FindPending(ctx context.Context, sentBeforeAt time.Time, tenantIDs []string) ([]*models.EnvelopeModel, error)
	FindByTxHashes(ctx context.Context, hashes, tenantIDs []string) ([]*models.EnvelopeModel, error)
	UpdateStatus(ctx context.Context, obj *models.EnvelopeModel, tenantIDs []string) error
}
