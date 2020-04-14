package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

//go:generate mockgen -source=load_envelope_by_id.go -destination=mocks/load_envelope_by_id.go -package=mocks

type LoadEnvelopeByID interface {
	Execute(ctx context.Context, tenantID string, envelopeID string) (*models.EnvelopeModel, error)
}

// RegisterContract is a use case to register a new contract
type loadEnvelopeByID struct {
	envelopeAgent store.EnvelopeAgent
}

// NewGetCatalog creates a new GetCatalog
func NewLoadEnvelopeByID(envelopeAgent store.EnvelopeAgent) LoadEnvelopeByID {
	return &loadEnvelopeByID{
		envelopeAgent: envelopeAgent,
	}
}

func (se *loadEnvelopeByID) Execute(ctx context.Context, tenantID, envelopeID string) (*models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	envelope, err := se.envelopeAgent.FindByFieldSet(ctx, map[string]string{
		"envelope_id": envelopeID,
		"tenant_id":   tenantID,
	})

	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"envelope.id": envelopeID,
				"tenant":      tenantID,
			}).
			Debugf("could not load envelope")
		if err == pg.ErrNoRows {
			return nil, errors.NotFoundError("envelope with id %v does not exist", envelopeID)
		}
		return nil, errors.StorageError(err.Error())
	}

	return envelope, nil
}
