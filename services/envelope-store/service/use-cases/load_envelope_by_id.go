package use_cases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

type LoadEnvelopeById interface {
	Execute(ctx context.Context, tenantId string, envelopeId string) (models.EnvelopeModel, error)
}

// RegisterContract is a use case to register a new contract
type loadEnvelopeById struct {
	envelopeAgent store.EnvelopeAgent
}

// NewGetCatalog creates a new GetCatalog
func NewLoadEnvelopeById(envelopeAgent store.EnvelopeAgent) LoadEnvelopeById {
	return &loadEnvelopeById{
		envelopeAgent: envelopeAgent,
	}
}

func (se *loadEnvelopeById) Execute(ctx context.Context, tenantId string, envelopeId string) (models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	envelope, err := se.envelopeAgent.FindByFieldSet(ctx, map[string]string{
		"envelope_id": envelopeId,
		"tenant_id": tenantId,
	})

	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"envelope.id": envelopeId,
				"tenant":   tenantId,
			}).
			Debugf("could not load envelope")
		if err == pg.ErrNoRows {
			return models.EnvelopeModel{}, errors.NotFoundError("envelope with id %v does not exist", envelopeId)
		}
		return models.EnvelopeModel{}, errors.StorageError(err.Error())
	}
	
	return envelope, nil
}
